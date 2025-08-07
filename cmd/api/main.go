package main

import (
	"fmt"
	"log"

	_ "sentiment-api/docs" // Import swagger docs
	"sentiment-api/internal/client"
	"sentiment-api/internal/config"
	"sentiment-api/internal/handler"
	"sentiment-api/internal/service"
	"sentiment-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Log.Level, cfg.Log.Format)
	logger.LogInfo("Starting Sentiment Analysis API", logrus.Fields{
		"port": cfg.Server.Port,
		"host": cfg.Server.Host,
	})

	// Validate required configuration
	if cfg.LLM.APIKey == "" {
		logger.LogError("LLM_API_KEY environment variable is required", nil)
		log.Fatal("LLM_API_KEY environment variable is required")
	}

	if cfg.LLM.URL == "" {
		logger.LogError("URL_CHAT_LLM_LLM environment variable is required", nil)
		log.Fatal("URL_CHAT_LLM_LLM environment variable is required")
	}

	// Initialize clients
	llmClient := client.NewLLMClient(cfg)

	// Initialize services
	sentimentService := service.NewSentimentService(llmClient)

	// Initialize handlers
	sentimentHandler := handler.NewSentimentHandler(sentimentService)

	// Setup router
	router := setupRouter(sentimentHandler)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.LogInfo("Server starting", logrus.Fields{
		"address": address,
	})

	if err := router.Run(address); err != nil {
		logger.LogError("Failed to start server", logrus.Fields{
			"error": err.Error(),
		})
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRouter configures and returns the Gin router
func setupRouter(sentimentHandler *handler.SentimentHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", sentimentHandler.HealthCheck)

	// Swagger documentation endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		sentiment := v1.Group("/sentiment")
		{
			sentiment.POST("/analyze", sentimentHandler.AnalyzeSentiment)
			sentiment.GET("/types", sentimentHandler.GetSentiments)
		}
	}

	return router
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
