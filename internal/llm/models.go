package llm

import "encoding/json"

type TogetherAiLLMRequestModel struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TogetherAiLLMResponseModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Prompt  []any  `json:"prompt"`
	Choices []struct {
		FinishReason string      `json:"finish_reason"`
		Seed         json.Number `json:"seed"`
		Logprobs     any         `json:"logprobs"`
		Index        int         `json:"index"`
		Message      struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []any  `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
		CachedTokens     int `json:"cached_tokens"`
	} `json:"usage"`
}
