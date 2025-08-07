package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger initializes the logger
func InitLogger(level, format string) {
	Log = logrus.New()

	// Set log level
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	Log.SetOutput(os.Stdout)
}

// LogErrorWithContext logs error with context
func LogErrorWithContext(err error, context string) {
	Log.WithFields(logrus.Fields{
		"context": context,
		"error":   err.Error(),
	}).Error("Error occurred")
}

// LogInfo logs info message
func LogInfo(message string, fields logrus.Fields) {
	if fields != nil {
		Log.WithFields(fields).Info(message)
	} else {
		Log.Info(message)
	}
}

// LogDebug logs debug message
func LogDebug(message string, fields logrus.Fields) {
	if fields != nil {
		Log.WithFields(fields).Debug(message)
	} else {
		Log.Debug(message)
	}
}

// LogError logs error message
func LogError(message string, fields logrus.Fields) {
	if fields != nil {
		Log.WithFields(fields).Error(message)
	} else {
		Log.Error(message)
	}
}

// LogWarn logs warning message
func LogWarn(message string, fields logrus.Fields) {
	if fields != nil {
		Log.WithFields(fields).Warn(message)
	} else {
		Log.Warn(message)
	}
}
