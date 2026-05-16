package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"auth-service/internal/config"
	"nexus-erp/shared/pkg/logger"
)

var (
	db     *pgxpool.Pool
	zapLog *zap.Logger
)

func main() {
	// Initialize logger
	zapLog, _ = zap.NewProduction()
	defer zapLog.Sync()

	// Load configuration
	cfg := config.Load()

	// Database connection
	dbURL := cfg.DatabaseURL()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	db, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		zapLog.Fatal("Cannot create database pool", zap.Error(err))
	}
	defer db.Close()

	err = db.Ping(ctx)
	if err != nil {
		zapLog.Fatal("Database ping failed", zap.Error(err))
	}

	zapLog.Info("Connected to PostgreSQL",
		zap.String("host", cfg.DBHost+":"+cfg.DBPort),
		zap.String("database", cfg.DBName),
	)

	// Gin router with custom logger middleware
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.GinMiddleware(zapLog))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "auth-service",
			"status":  "healthy",
		})
	})

	zapLog.Info("Auth Service starting", zap.String("port", cfg.ServerPort))
	r.Run(":" + cfg.ServerPort)
}