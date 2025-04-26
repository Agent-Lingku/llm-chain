package gen

import (
	"encoding/json"
	"fmt"
	"learn/internal/interfaces"
	"time"

	"resty.dev/v3"
)

// RemoteLargeModelClient 远程大模型客户端 (使用Resty实现)
type RemoteLargeModelClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

// NewRemoteLargeModelClient 创建新的远程大模型客户端
func NewRemoteLargeModelClient(baseURL, apiKey string) *RemoteLargeModelClient {
	client := resty.New()
	client.
		SetTimeout(300*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5*time.Second).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Accept", "application/json")

	return &RemoteLargeModelClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// 实现 IRemoteLLM 接口
func (c *RemoteLargeModelClient) Generate(input string) (string, error) {
	requestBody := map[string]interface{}{
		"model": "gpt-4-turbo",
		"messages": []interfaces.Message{
			{Role: "user", Content: input},
		},
		"temperature": 0.7,
	}

	resp, err := c.client.R().
		SetBody(requestBody).
		Post(c.baseURL + "/v1/chat/completions")
	if err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("API异常响应: %s\n%s", resp.Status(), string(resp.Bytes()))
	}

	var result struct {
		Choices []struct {
			Message interfaces.Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(resp.Bytes(), &result); err != nil {
		return "", fmt.Errorf("JSON解析失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API返回空结果")
	}
	return result.Choices[0].Message.Content, nil
}

func (c *RemoteLargeModelClient) GetModelParams() map[string]string {
	return map[string]string{
		"model":       "gpt-4-turbo",
		"temperature": "0.7",
	}
}

func (c *RemoteLargeModelClient) ParseResponse(response string) (string, error) {
	// 使用gjson解析
	content := interfaces.ParseJSONField(response, "choices.0.message.content")
	if content == "" {
		return "", fmt.Errorf("无法解析响应内容")
	}
	return content, nil
}

// Message 定义API消息结构（移到interfaces包）
// 注意：需确保interfaces.Message类型与这里定义一致
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
