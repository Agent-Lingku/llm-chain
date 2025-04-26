package gen

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"learn/internal/config"

	"github.com/tidwall/gjson"
	"resty.dev/v3"
)

// 定义模型列表结构
type ModelList struct {
	Models []Models `json:"models"`
}

// 定义模型详情结构
type Details struct {
	Format            string `json:"format"`
	Family            string `json:"family"`
	Families          any    `json:"families"`
	ParameterSize     string `json:"parameter_size"`
	QuantizationLevel string `json:"quantization_level"`
}

// 定义单个模型结构
type Models struct {
	Name       string  `json:"name"`
	ModifiedAt string  `json:"modified_at"`
	Size       int64   `json:"size"`
	Digest     string  `json:"digest"`
	Details    Details `json:"details"`
}

// 定义模型信息结构
type ModelInfo struct {
	ModelFile  string         `json:"modelfile"`
	Parameters string         `json:"parameters"`
	Template   string         `json:"template"`
	Details    map[string]any `json:"details"`
	ModelInfo  map[string]any `json:"model_info"`
}

// 定义本地大语言模型接口
type ILocalLLM interface {
	ModelList() ([]Models, error)
	GetReply(resp *resty.Response) string
	Client(requestBody map[string]any) (*resty.Response, error)
}

// 定义本地大语言模型结构
type LocalLLM struct {
	client *resty.Client
}

// 创建本地大语言模型客户端实例
func NewLocalLargeModelClient() ILocalLLM {
	client := resty.New()
	client.SetBaseURL(config.OllamaUrl)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	client.SetTimeout(120 * time.Second)
	return &LocalLLM{
		client: client,
	}
}

// 获取模型列表
func (llm *LocalLLM) ModelList() ([]Models, error) {
	resp, err := llm.client.R().Get("/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to get model list: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get model list: HTTP %d", resp.StatusCode())
	}

	var modelList ModelList
	if err := json.Unmarshal(resp.Bytes(), &modelList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal model list: %w", err)
	}

	return modelList.Models, nil
}

func (llm *LocalLLM) GetReply(resp *resty.Response) string {
	if resp.IsError() {
		return ""
	}
	return gjson.GetBytes(resp.Bytes(), "message.content").String()
}

// 生成聊天响应
func (llm *LocalLLM) Client(requestBody map[string]any) (*resty.Response, error) {
	return llm.client.R().SetBody(requestBody).Post(config.OllamaPrefix)
}
