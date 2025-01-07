package handler

import (
	"context"
	"github.com/optimism-java/interopbankend/internal/schema"
	"github.com/optimism-java/interopbankend/internal/svc"
	"github.com/optimism-java/interopbankend/pkg/log"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

func LatestBlackNumber(ctx *svc.ServiceContext) {
	for {
		latest, err := ctx.L2RPC.BlockNumber(context.Background())
		if err != nil {
			log.Errorf("[Handler.LatestBlackNumber] Syncing block by number error: %s\n", errors.WithStack(err))
			time.Sleep(12 * time.Second)
			continue
		}

		ctx.LatestBlockNumber = cast.ToInt64(latest)

		time.Sleep(12 * time.Second)
	}
}

func SyncFinalBlock(ctx *svc.ServiceContext) {
	for {
		var block map[string]interface{}
		err := ctx.L2RPC.Client().CallContext(context.Background(), &block, "eth_getBlockByNumber", "finalized", true)
		if err != nil {
			log.Fatal(err.Error())
		}
		finalizedBlockNumber := cast.ToInt64(block["number"].(string))

		result := ctx.DB.Model(&schema.SyncBlock{}).
			Where("block_number <= ? AND block_state != ?", finalizedBlockNumber, "finalized").
			Update("block_state", "finalized")

		if result.Error != nil {
			log.Errorf("[Handler.SyncFinalBlock] Batch update blocks error: %s\n", errors.WithStack(result.Error))
		}

		time.Sleep(24 * time.Second)
	}
}
