package interfaces

import "github.com/tidwall/gjson"

// IRemoteLLM defines the interface for remote LLM operations
type IRemoteLLM interface {
    Generate(input string) (string, error)
    GetModelParams() map[string]string
    ParseResponse(response string) (string, error) // Uses gjson for JSON parsing
}

// ILocalLLM defines the interface for local LLM operations
type ILocalLLM interface {
    Generate(input string) (string, error)
    GetModelParams() map[string]string
    LoadLocalModel(modelPath string) error
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// Common JSON parsing helper using gjson
func ParseJSONField(jsonStr, field string) string {
    return gjson.Get(jsonStr, field).String()
}
