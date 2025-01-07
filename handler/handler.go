package handler

import "github.com/optimism-java/interopbackend/internal/svc"

func Run(ctx *svc.ServiceContext) {
	// query last block number
	go LatestBlackNumber(ctx)
	// sync blocks
	go SyncBlock(ctx)
	// sync block finalized state
	go SyncFinalBlock(ctx)
	// sync events
	go SyncEvent(ctx)
}
