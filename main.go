package main

import (
	"context"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/optimism-java/interopbackend/docs" // import swagger docs
	"github.com/optimism-java/interopbackend/handler"
	"github.com/optimism-java/interopbackend/internal/api"
	"github.com/optimism-java/interopbackend/internal/svc"
	"github.com/optimism-java/interopbackend/internal/types"
	"github.com/optimism-java/interopbackend/migration/migrate"
	"github.com/optimism-java/interopbackend/pkg/log"
)

// @title Interop Backend API
// @version 1.0
// @description Interop Backend Service API Documentation
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	cfg := types.GetConfig()
	log.Init(cfg.LogLevel, cfg.LogFormat)
	log.Infof("config: %v\n", cfg)
	sCtx := svc.NewServiceContext(ctx, cfg)
	migrate.Migrate(sCtx.DB)
	handler.Run(sCtx)
	log.Info("listener running...\n")
	router := gin.Default()
	apiHandler := api.NewAPIHandler(sCtx.Config, sCtx.DB)

	// API Routes
	router.GET("/blocks", apiHandler.GetSyncBlocks)
	router.GET("/blocks/:blockNumber/executingMessage", apiHandler.GetExecutingMessageByBlockNumber)
	router.GET("/blocks/sentMessage/:hash", apiHandler.GetSentMessageByHash)
	router.GET("/blocks/relayedMessage/:hash", apiHandler.GetRelayedMessageByHash)

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := router.Run(":" + cfg.APIPort)
	if err != nil {
		log.Errorf("start error %s", err)
		return
	}
}
