package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"sentiment-api/internal/config"
	"sentiment-api/internal/model"
	"sentiment-api/pkg/logger"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// LLMClient handles communication with LLM API
type LLMClient struct {
	config *config.Config
	client *resty.Client
}

// NewLLMClient creates a new LLM client
func NewLLMClient(cfg *config.Config) *LLMClient {
	client := resty.New()
	client.SetTimeout(60 * time.Second)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("x-api-key", cfg.LLM.APIKey)

	return &LLMClient{
		config: cfg,
		client: client,
	}
}

// CallTelkomAI makes a call to Telkom AI API
func (c *LLMClient) CallTelkomAI(messages []model.LLMMessage, modelName string, maxTokens int, temperature float64) (interface{}, error) {
	logger.LogDebug("Making API call to LLM", logrus.Fields{
		"model":       modelName,
		"messages":    len(messages),
		"max_tokens":  maxTokens,
		"temperature": temperature,
	})

	request := model.LLMRequest{
		Model:       modelName,
		Messages:    messages,
		Stream:      false,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	var response model.LLMResponse
	resp, err := c.client.R().
		SetBody(request).
		SetResult(&response).
		Post(c.config.LLM.URL)

	if err != nil {
		logger.LogErrorWithContext(err, "HTTP request error in LLM call")
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		errorMsg := fmt.Sprintf("response error %d: %s", resp.StatusCode(), resp.String())
		logger.LogError("HTTP status error", logrus.Fields{
			"status_code": resp.StatusCode(),
			"response":    resp.String(),
		})
		return nil, errors.New(errorMsg)
	}

	if len(response.Choices) == 0 {
		logger.LogError("No choices in LLM response", nil)
		return nil, errors.New("no choices in response")
	}

	content := response.Choices[0].Message.Content

	logger.LogDebug("LLM API call successful", logrus.Fields{
		"content_length": len(content),
	})

	// Try to parse JSON response
	var parsedContent interface{}
	if err := json.Unmarshal([]byte(content), &parsedContent); err != nil {
		logger.LogWarn("Content is not valid JSON, returning as string", logrus.Fields{
			"error":   err.Error(),
			"content": content,
		})
		return content, nil
	}

	logger.LogDebug("JSON parsing successful", nil)
	return parsedContent, nil
}

// AnalyzeSentiment performs sentiment analysis using LLM
func (c *LLMClient) AnalyzeSentiment(textPertanyaan, textJawaban string) (string, error) {
	systemPrompt := `Anda adalah sistem analisis sentimen yang sangat akurat. Tugas Anda adalah menganalisis sentimen dari jawaban terhadap pertanyaan yang diberikan.

Berdasarkan konteks pertanyaan dan jawaban, tentukan sentimen jawaban tersebut:
- Positif: Jawaban menunjukkan emosi atau pandangan yang baik, puas, senang, atau mendukung
- Negatif: Jawaban menunjukkan emosi atau pandangan yang buruk, tidak puas, kecewa, atau menolak  
- Netral: Jawaban objektif, tidak menunjukkan emosi khusus, atau seimbang

Respons Anda harus dalam format JSON yang valid:
{"sentiment": "Positif"} atau {"sentiment": "Negatif"} atau {"sentiment": "Netral"}

Hanya gunakan kata: Positif, Negatif, atau Netral.`

	userPrompt := fmt.Sprintf(`Pertanyaan: %s

Jawaban: %s

Analisis sentimen jawaban tersebut berdasarkan konteks pertanyaan.`, textPertanyaan, textJawaban)

	messages := []model.LLMMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	result, err := c.CallTelkomAI(messages, "telkom-ai-instruct", 100, 0.0)
	if err != nil {
		return "", err
	}

	// Parse the result to extract sentiment
	sentiment, err := c.extractSentimentFromResult(result)
	if err != nil {
		logger.LogError("Failed to extract sentiment from LLM response", logrus.Fields{
			"result": result,
			"error":  err.Error(),
		})
		return "Netral", nil // Default to Netral if parsing fails
	}

	return sentiment, nil
}

// AnalyzeSentimentWithReasoning performs sentiment analysis with reasoning explanation using LLM
func (c *LLMClient) AnalyzeSentimentWithReasoning(textPertanyaan, textJawaban string) (string, *string, error) {
	systemPrompt := `Anda adalah sistem analisis sentimen yang sangat akurat dan dapat memberikan penjelasan. Tugas Anda adalah menganalisis sentimen dari jawaban terhadap pertanyaan yang diberikan, beserta alasan analisis tersebut.

Berdasarkan konteks pertanyaan dan jawaban, tentukan sentimen jawaban tersebut:
- Positif: Jawaban menunjukkan emosi atau pandangan yang baik, puas, senang, atau mendukung
- Negatif: Jawaban menunjukkan emosi atau pandangan yang buruk, tidak puas, kecewa, atau menolak  
- Netral: Jawaban objektif, tidak menunjukkan emosi khusus, atau seimbang

Respons Anda harus dalam format JSON yang valid dengan penjelasan:
{
  "sentiment": "Positif",
  "reasoning": "Penjelasan mengapa sentimen ini dipilih, kata-kata kunci yang mendukung, dan konteks yang relevan"
}

Hanya gunakan kata: Positif, Negatif, atau Netral untuk sentiment.
Berikan penjelasan yang jelas dan informatif dalam bahasa Indonesia untuk reasoning.`

	userPrompt := fmt.Sprintf(`Pertanyaan: %s

Jawaban: %s

Analisis sentimen jawaban tersebut berdasarkan konteks pertanyaan dan berikan penjelasan lengkap.`, textPertanyaan, textJawaban)

	messages := []model.LLMMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	result, err := c.CallTelkomAI(messages, "telkom-ai-instruct", 300, 0.1)
	if err != nil {
		return "", nil, err
	}

	// Parse the result to extract sentiment and reasoning
	sentiment, reasoning, err := c.extractSentimentAndReasoningFromResult(result)
	if err != nil {
		logger.LogError("Failed to extract sentiment and reasoning from LLM response", logrus.Fields{
			"result": result,
			"error":  err.Error(),
		})
		// Return basic sentiment without reasoning on parsing failure
		basicSentiment, basicErr := c.extractSentimentFromResult(result)
		if basicErr != nil {
			return "Netral", nil, nil // Default to Netral if everything fails
		}
		return basicSentiment, nil, nil
	}

	return sentiment, reasoning, nil
}

// extractSentimentAndReasoningFromResult extracts both sentiment and reasoning from LLM result
func (c *LLMClient) extractSentimentAndReasoningFromResult(result interface{}) (string, *string, error) {
	// If result is a map (parsed JSON)
	if resultMap, ok := result.(map[string]interface{}); ok {
		sentiment, sentimentExists := resultMap["sentiment"]
		reasoning, reasoningExists := resultMap["reasoning"]

		if sentimentExists {
			sentimentStr, sentimentOk := sentiment.(string)
			if sentimentOk {
				normalizedSentiment := c.normalizeSentiment(sentimentStr)

				if reasoningExists {
					if reasoningStr, reasoningOk := reasoning.(string); reasoningOk && reasoningStr != "" {
						return normalizedSentiment, &reasoningStr, nil
					}
				}
				return normalizedSentiment, nil, nil
			}
		}
	}

	// If result is a string, try to parse it as JSON
	if resultStr, ok := result.(string); ok {
		var sentimentResult map[string]interface{}
		if err := json.Unmarshal([]byte(resultStr), &sentimentResult); err == nil {
			sentiment, sentimentExists := sentimentResult["sentiment"]
			reasoning, reasoningExists := sentimentResult["reasoning"]

			if sentimentExists {
				if sentimentStr, ok := sentiment.(string); ok {
					normalizedSentiment := c.normalizeSentiment(sentimentStr)

					if reasoningExists {
						if reasoningStr, ok := reasoning.(string); ok && reasoningStr != "" {
							return normalizedSentiment, &reasoningStr, nil
						}
					}
					return normalizedSentiment, nil, nil
				}
			}
		}

		// If not proper JSON, try to extract sentiment only
		extractedSentiment := c.extractSentimentFromString(resultStr)
		return extractedSentiment, nil, nil
	}

	return "", nil, errors.New("unable to extract sentiment and reasoning from result")
}

// extractSentimentFromResult extracts sentiment from LLM result
func (c *LLMClient) extractSentimentFromResult(result interface{}) (string, error) {
	// If result is a map (parsed JSON)
	if resultMap, ok := result.(map[string]interface{}); ok {
		if sentiment, exists := resultMap["sentiment"]; exists {
			if sentimentStr, ok := sentiment.(string); ok {
				return c.normalizeSentiment(sentimentStr), nil
			}
		}
	}

	// If result is a string, try to parse it as JSON
	if resultStr, ok := result.(string); ok {
		var sentimentResult map[string]interface{}
		if err := json.Unmarshal([]byte(resultStr), &sentimentResult); err == nil {
			if sentiment, exists := sentimentResult["sentiment"]; exists {
				if sentimentStr, ok := sentiment.(string); ok {
					return c.normalizeSentiment(sentimentStr), nil
				}
			}
		}

		// If not JSON, try to extract sentiment directly from string
		return c.extractSentimentFromString(resultStr), nil
	}

	return "", errors.New("unable to extract sentiment from result")
}

// extractSentimentFromString extracts sentiment from string response
func (c *LLMClient) extractSentimentFromString(text string) string {
	text = fmt.Sprintf("%s", text) // Convert to lowercase for comparison

	if contains(text, "positif") {
		return "Positif"
	} else if contains(text, "negatif") {
		return "Negatif"
	} else if contains(text, "netral") {
		return "Netral"
	}

	// Default to Netral if no clear sentiment found
	return "Netral"
}

// normalizeSentiment normalizes sentiment values
func (c *LLMClient) normalizeSentiment(sentiment string) string {
	switch sentiment {
	case "Positif", "positif", "POSITIF", "Positive", "positive", "POSITIVE":
		return "Positif"
	case "Negatif", "negatif", "NEGATIF", "Negative", "negative", "NEGATIVE":
		return "Negatif"
	case "Netral", "netral", "NETRAL", "Neutral", "neutral", "NEUTRAL":
		return "Netral"
	default:
		return "Netral"
	}
}

// contains checks if a string contains a substring (case insensitive)
func contains(text, substr string) bool {
	return len(text) >= len(substr) &&
		(text == substr ||
			len(text) > len(substr) &&
				(text[0:len(substr)] == substr ||
					text[len(text)-len(substr):] == substr ||
					findSubstring(text, substr)))
}

// findSubstring finds substring in text
func findSubstring(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
