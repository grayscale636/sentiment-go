package model

// SentimentRequest represents the input for sentiment analysis
type SentimentRequest struct {
	TextPertanyaan string `json:"text_pertanyaan" binding:"required" example:"Bagaimana pendapat Anda tentang layanan kami?" description:"The question or prompt text"`
	TextJawaban    string `json:"text_jawaban" binding:"required" example:"Layanan Anda sangat memuaskan dan responsif" description:"The answer or response text to be analyzed"`
	Reasoning      *bool  `json:"reasoning,omitempty" example:"true" description:"Optional: Request reasoning explanation from LLM (default: false)"`
}

// SentimentResponse represents the output of sentiment analysis
type SentimentResponse struct {
	Sentiment string  `json:"sentiment" example:"Positif" enum:"Positif,Negatif,Netral" description:"The analyzed sentiment: Positif (positive), Negatif (negative), or Netral (neutral)"`
	Reasoning *string `json:"reasoning,omitempty" example:"Teks menunjukkan kepuasan pelanggan dengan kata-kata positif seperti 'memuaskan' dan 'responsif'" description:"Optional: LLM reasoning explanation for the sentiment analysis"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message" example:"text_pertanyaan and text_jawaban are required"`
}

// APIResponse represents general API response
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// LLMMessage represents message structure for LLM API
type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMRequest represents request to LLM API
type LLMRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	Stream      bool         `json:"stream"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`
}

// LLMChoice represents choice in LLM response
type LLMChoice struct {
	Message LLMMessage `json:"message"`
	Index   int        `json:"index"`
}

// LLMResponse represents response from LLM API
type LLMResponse struct {
	Choices []LLMChoice `json:"choices"`
	Model   string      `json:"model"`
	Usage   interface{} `json:"usage"`
}
