package util

import (
	"regexp"
	"strings"
)

type PromptType struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func AppendSystemPrompt(prompt string) PromptType {
	return PromptType{Role: "system", Content: prompt}
}

func AppendUserPrompt(prompt string) PromptType {
	return PromptType{Role: "user", Content: prompt}
}

func AppendAssistantPrompt(prompt string) PromptType {
	return PromptType{Role: "assistant", Content: prompt}
}

func ExtractCodeBlocks(markdown string) []string {
	codeBlocks := []string{}
	regex := regexp.MustCompile("```(?:[\\w-]+)?\\n([\\s\\S]*?)```")
	matches := regex.FindAllStringSubmatch(markdown, -1)

	for _, match := range matches {
		code := strings.TrimSpace(match[1])
		codeBlocks = append(codeBlocks, code)
	}

	return codeBlocks
}

func ForEach[T any](slice []T, operation func(T) error) error {
	for _, item := range slice {
		if err := operation(item); err != nil {
			return err
		}
	}
	return nil
}
