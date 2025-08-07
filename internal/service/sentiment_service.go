package service

import (
	"errors"
	"strings"

	"sentiment-api/internal/client"
	"sentiment-api/internal/model"
	"sentiment-api/pkg/logger"

	"github.com/sirupsen/logrus"
)

// SentimentService handles sentiment analysis business logic
type SentimentService struct {
	llmClient *client.LLMClient
}

// NewSentimentService creates a new sentiment service
func NewSentimentService(llmClient *client.LLMClient) *SentimentService {
	return &SentimentService{
		llmClient: llmClient,
	}
}

// AnalyzeSentiment analyzes sentiment of the given text pair
func (s *SentimentService) AnalyzeSentiment(req *model.SentimentRequest) (*model.SentimentResponse, error) {
	logger.LogInfo("Starting sentiment analysis", logrus.Fields{
		"text_pertanyaan_length": len(req.TextPertanyaan),
		"text_jawaban_length":    len(req.TextJawaban),
		"reasoning_requested":    req.Reasoning != nil && *req.Reasoning,
	})

	// Validate input
	if err := s.validateRequest(req); err != nil {
		logger.LogError("Request validation failed", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	// Check if reasoning is requested
	requestReasoning := req.Reasoning != nil && *req.Reasoning

	// Perform sentiment analysis using LLM
	var sentiment string
	var reasoning *string
	var err error

	if requestReasoning {
		sentiment, reasoning, err = s.llmClient.AnalyzeSentimentWithReasoning(req.TextPertanyaan, req.TextJawaban)
	} else {
		sentiment, err = s.llmClient.AnalyzeSentiment(req.TextPertanyaan, req.TextJawaban)
	}

	if err != nil {
		logger.LogError("Failed to analyze sentiment", logrus.Fields{
			"error": err.Error(),
		})
		return nil, err
	}

	logger.LogInfo("Sentiment analysis completed", logrus.Fields{
		"sentiment":         sentiment,
		"reasoning_present": reasoning != nil,
	})

	response := &model.SentimentResponse{
		Sentiment: sentiment,
		Reasoning: reasoning,
	}

	return response, nil
}

// validateRequest validates the sentiment analysis request
func (s *SentimentService) validateRequest(req *model.SentimentRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if strings.TrimSpace(req.TextPertanyaan) == "" {
		return errors.New("text_pertanyaan cannot be empty")
	}

	if strings.TrimSpace(req.TextJawaban) == "" {
		return errors.New("text_jawaban cannot be empty")
	}

	// Optional: Add length validation
	if len(req.TextPertanyaan) > 1000 {
		return errors.New("text_pertanyaan exceeds maximum length of 1000 characters")
	}

	if len(req.TextJawaban) > 2000 {
		return errors.New("text_jawaban exceeds maximum length of 2000 characters")
	}

	return nil
}

// GetSupportedSentiments returns list of supported sentiment values
func (s *SentimentService) GetSupportedSentiments() []string {
	return []string{"Positif", "Negatif", "Netral"}
}
