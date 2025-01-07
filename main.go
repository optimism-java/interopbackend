package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/optimism-java/interopbankend/handler"
	"github.com/optimism-java/interopbankend/internal/api"
	"github.com/optimism-java/interopbankend/internal/svc"
	"github.com/optimism-java/interopbankend/internal/types"
	"github.com/optimism-java/interopbankend/migration/migrate"
	"github.com/optimism-java/interopbankend/pkg/log"
)

func main() {
	ctx := context.Background()
	cfg := types.GetConfig()
	log.Init(cfg.LogLevel, cfg.LogFormat)
	log.Infof("config: %v\n", cfg)
	sCtx := svc.NewServiceContext(ctx, cfg)
	migrate.Migrate(sCtx.DB)
	handler.Run(sCtx)
	log.Info("listener running...\n")
	log.Info("listener running...\n")
	router := gin.Default()
	apiHandler := api.NewAPIHandler(sCtx.Config, sCtx.DB)

	router.GET("/blocks", apiHandler.GetSyncBlocks)
	router.GET("/blocks/:blockNumber/executingMessage", apiHandler.GetExecutingMessageByBlockNumber)
	router.GET("/sentMessage/:hash", apiHandler.GetSentMessageByHash)
	router.GET("/relayedMessage/:hash", apiHandler.GetRelayedMessageByHash)

	err := router.Run(":" + cfg.APIPort)
	if err != nil {
		log.Errorf("start error %s", err)
		return
	}
}