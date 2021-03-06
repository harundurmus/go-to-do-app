package main

import (
	"github.com/harundurmus/go-to-do-app/config"
	_ "github.com/harundurmus/go-to-do-app/docs"
	"github.com/harundurmus/go-to-do-app/internal/client"
	custommiddleware "github.com/harundurmus/go-to-do-app/internal/middleware"
	"github.com/harundurmus/go-to-do-app/internal/todoapp"
	_ "github.com/harundurmus/go-to-do-app/internal/todoapp"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strings"
)

func main() {

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	conf, err := config.New(".config", appEnv)
	if err != nil {
		panic(err)
	}
	logger := buildLogger(conf.LogLevel)

	server := NewServer(conf)
	connectMongoDb := client.ConnectMongoDb(conf.MongoDB)
	_ = custommiddleware.AuthMiddleware{
		SecretKey: conf.SecretKey,
		Aud:       conf.Aud,
		Iss:       conf.Iss,
	}
	repository := todoapp.NewRepository(connectMongoDb, conf.MongoDB)
	service := todoapp.NewService(repository, logger)
	handler := todoapp.NewHandler(logger, service)
	server.e.GET("/swagger/*", echoSwagger.WrapHandler)
	server.e.POST("/init", handler.InitializeDatabase)
	server.e.POST("/create-task", handler.CreateTaskData)
	_ = server.Start()
}

func buildLogger(logLevel string) *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(getLogLevel(logLevel))
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func getLogLevel(level string) zapcore.Level {
	switch levelFromConfig := strings.TrimSpace(level); {
	case strings.EqualFold(levelFromConfig, "debug"):
		return zapcore.DebugLevel
	case strings.EqualFold(levelFromConfig, "error"):
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
