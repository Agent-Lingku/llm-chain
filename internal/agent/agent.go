package agent

import (
	"encoding/json"
	"fmt"
	"learn/internal/config"
	"learn/internal/model"
	"learn/internal/util"
	"strings"
	"time"

	"log"

	"github.com/tidwall/gjson"
	"resty.dev/v3"
)

// Agent 代表一个代理，用于执行特定的任务
type Agent struct {
	client *resty.Client
	config AConfig
}

// State 状态
type State int

const (
	StatePending State = iota
	StateRunning
	StateCompleted
	StateFailed
)

// AConfig 代理配置
type AConfig struct {
	TaskID       string
	AgentName    string
	Model        string
	Role         Role
	SystemPrompt string
	UserPrompt   string
	Status       model.Status
	Context      []map[string]string
	EnableSearch bool
	CurrentState State
}

// Option 定义 with 选项函数类型
type Option func(*AConfig)

// AdapterFunc 适配器承载不同模型对接
type AdapterFunc func(body map[string]any) (*resty.Response, error)

// WithTaskID 设置 TaskID
func WithTaskID(taskID string) Option {
	return func(cfg *AConfig) {
		cfg.TaskID = taskID
	}
}

// WithAgentName 设置 AgentName
func WithAgentName(agentName string) Option {
	return func(cfg *AConfig) {
		cfg.AgentName = agentName
	}
}

// WithModel 设置 Model
func WithModel(model string) Option {
	return func(cfg *AConfig) {
		cfg.Model = model
	}
}

// WithRole 设置 Role
func WithRole(role Role) Option {
	return func(cfg *AConfig) {
		cfg.Role = role
	}
}

// WithUserPrompt 设置 UserPrompt
func WithUserPrompt(userPrompt string) Option {
	return func(cfg *AConfig) {
		cfg.UserPrompt = userPrompt
	}
}

// WithContext 设置 Context
func WithContext(context []map[string]string) Option {
	return func(cfg *AConfig) {
		cfg.Context = context
	}
}

// WithEnableSearch 设置 EnableSearch
func WithEnableSearch(enableSearch bool) Option {
	return func(cfg *AConfig) {
		cfg.EnableSearch = enableSearch
	}
}

// WithStatus 设置 Status
func WithStatus(status model.Status) Option {
	return func(cfg *AConfig) {
		cfg.Status = status
	}
}

// NewAgent 创建一个新的Agent
func NewAgent(opts ...Option) *Agent {
	client := resty.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	client.SetBaseURL(cfg.ApiBaseUrl)

	agent := &Agent{
		client: client,
		config: AConfig{
			Model:  "qwen-max",
			Status: model.StatusPending,
		},
	}

	for _, opt := range opts {
		opt(&agent.config)
	}

	agent.config.SystemPrompt = GetAgentPrompt(agent.config.Role)

	return agent
}

func (a *Agent) EchoRoleInfo() (string, string) {
	return a.config.AgentName, a.config.SystemPrompt
}

func (a *Agent) IsFinished() bool {
	return a.config.Status == model.StatusCompleted
}

func (a *Agent) HasContext() bool {
	return len(a.config.Context) > 0
}

func (a *Agent) GetStatus() model.Status {
	return a.config.Status
}

func (a *Agent) setStatus(status model.Status) {
	a.config.Status = status
}

func (a *Agent) extractContent(res *resty.Response) string {
	if strings.Contains(res.Request.URL, config.OllamaUrl) {
		return gjson.GetBytes(res.Bytes(), "message.content").String()
	}
	return gjson.GetBytes(res.Bytes(), "choices.0.message.content").String()
}

func (a *Agent) checkError(res *resty.Response) (string, error) {
	if res.IsError() {
		a.setStatus(model.StatusFailed)
		bodyString := string(res.Bytes())
		log.Printf("请求失败: %s, Body: %s\n", res.Status(), bodyString)
		return "", fmt.Errorf("request failed: %s", res.Status())
	}

	content := a.extractContent(res)
	if content == "" {
		a.setStatus(model.StatusFailed)
		return "", fmt.Errorf("empty content from API")
	}

	a.setStatus(model.StatusCompleted)
	return content, nil
}

// parseRequestBody 构造请求体
func (a *Agent) parseRequestBody(more ...util.PromptType) map[string]any {
	messages := []util.PromptType{
		util.AppendSystemPrompt(a.config.SystemPrompt),
		util.AppendUserPrompt(a.config.UserPrompt),
	}

	messages = append(messages, more...)

	for _, msg := range a.config.Context {
		if msg["role"] != "user" {
			continue
		}
		messages = append(messages, util.PromptType{
			Role:    msg["role"],
			Content: msg["content"],
		})
	}

	fmt.Println(messages)

	return map[string]any{
		"messages":      messages,
		"stream":        false,
		"model":         a.config.Model,
		"enable_search": a.config.EnableSearch,
	}
}

// RateLimiter 速率限制器
type RateLimiter struct {
	bucketSize int
	refillRate float64
	tokens     float64
	lastRefill time.Time
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(bucketSize int, refillRate float64) *RateLimiter {
	return &RateLimiter{
		bucketSize: bucketSize,
		refillRate: refillRate,
		tokens:     float64(bucketSize),
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow() bool {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > float64(rl.bucketSize) {
		rl.tokens = float64(rl.bucketSize)
	}
	rl.lastRefill = now

	if rl.tokens >= 1 {
		rl.tokens--
		return true
	}
	return false
}

// ToolCall 工具调用信息
type ToolCall struct {
	ToolName string                 `json:"tool_name"`
	Params   map[string]interface{} `json:"params"`
}

// ExecuteTask 执行任务并发送请求
func (a *Agent) ExecuteTask(callBack AdapterFunc, more ...util.PromptType) ([]ToolCall, string, error) {
	var res *resty.Response
	var err error
	var toolCalls []ToolCall

	// 实现速率限制
	rateLimiter := NewRateLimiter(10, 1.0) // 每秒最多10个请求

	// 实现指数退避重试机制
	retryCount := 3
	retryInterval := 1 * time.Second

	for i := 0; i < retryCount; i++ {
		if !rateLimiter.Allow() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		reqBody := a.parseRequestBody(more...)

		// 添加工具调用信息
		reqBody["tool_calls"] = []ToolCall{
			{
				ToolName: "code_editor",
				Params: map[string]interface{}{
					"action":  "edit",
					"file":    "main.go",
					"line":    10,
					"content": "fmt.Println(\"Hello, World!\")",
				},
			},
		}

		// call interface
		res, err = callBack(reqBody)

		if err == nil {
			break
		}

		time.Sleep(retryInterval)
		retryInterval *= 2 // 指数退避
	}

	if err != nil {
		return toolCalls, "", fmt.Errorf("failed to send request after %d retries: %w", retryCount, err)
	}

	// 解析工具调用信息
	toolCallsJSON := gjson.GetBytes(res.Bytes(), "tool_calls").String()
	if err := json.Unmarshal([]byte(toolCallsJSON), &toolCalls); err != nil {
		return toolCalls, "", fmt.Errorf("failed to parse tool calls: %w", err)
	}

	content, err := a.checkError(res)
	return toolCalls, content, err
}
