package message

type ChatResponse struct {
	Messages []Message  `json:"messages"`
	Usage    TokenUsage `json:"usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TokenUsage struct {
	TotalTokens  int `json:"total_tokens,omitempty"`
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type MetaData struct {
	TenantID  string
	RequestID string
	UserAgent string
	Feature   string
}
