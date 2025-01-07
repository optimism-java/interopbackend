package handler

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/optimism-java/interopbackend/internal/schema"
	"github.com/optimism-java/interopbackend/internal/svc"
	"github.com/optimism-java/interopbackend/pkg/log"
	"github.com/optimism-java/interopbackend/pkg/rpc"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func SyncBlock(ctx *svc.ServiceContext) {
	time.Sleep(10 * time.Second)
	var syncedBlock schema.SyncBlock
	err := ctx.DB.Where("status = ? or status = ? ", schema.BlockValid, schema.BlockPending).Order("block_number desc").First(&syncedBlock).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}
	switch err {
	case gorm.ErrRecordNotFound:
		ctx.SyncedBlockNumber = ctx.Config.FromBlockNumber
		ctx.SyncedBlockHash = common.HexToHash(ctx.Config.FromBlockHash)
	default:
		ctx.SyncedBlockNumber = syncedBlock.BlockNumber
		ctx.SyncedBlockHash = common.HexToHash(syncedBlock.BlockHash)
	}

	for {
		syncingBlockNumber := ctx.SyncedBlockNumber + 1
		log.Infof("[Handler.SyncBlock] Try to sync block number: %d\n", syncingBlockNumber)

		if syncingBlockNumber > ctx.LatestBlockNumber {
			time.Sleep(3 * time.Second)
			continue
		}

		// block, err := ctx.L2RPC.BlockByNumber(context.Background(), big.NewInt(syncingBlockNumber))
		blockJSON, err := rpc.HTTPPostJSON("", ctx.Config.L2RPCUrl, "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\""+fmt.Sprintf("0x%X", syncingBlockNumber)+"\", true],\"id\":1}")
		if err != nil {
			log.Errorf("[Handler.SyncBlock] Syncing block by number error: %s\n", errors.WithStack(err))
			time.Sleep(3 * time.Second)
			continue
		}
		block := rpc.ParseJSONBlock(string(blockJSON))
		log.Infof("[Handler.SyncBlock] Syncing block number: %d, hash: %v, parent hash: %v \n", block.Number(), block.Hash(), block.ParentHash())

		if common.HexToHash(block.ParentHash()) != ctx.SyncedBlockHash {
			log.Errorf("[Handler.SyncBlock] ParentHash of the block being synchronized is inconsistent: %s \n", ctx.SyncedBlockHash)
			rollbackBlock(ctx)
			continue
		}

		/* Create SyncBlock start */
		err = ctx.DB.Create(&schema.SyncBlock{
			Miner:        block.Result.Miner,
			Blockchain:   ctx.Config.Blockchain,
			BlockchainID: ctx.Config.BlockChainID,
			BlockTime:    block.Timestamp(),
			BlockNumber:  block.Number(),
			BlockHash:    block.Hash(),
			TxCount:      int64(len(block.Result.Transactions)),
			EventCount:   0,
			ParentHash:   block.ParentHash(),
			Status:       schema.BlockPending,
			BlockState:   schema.BlockStateLatest,
		}).Error
		if err != nil {
			log.Errorf("[Handler.SyncBlock] DB Create SyncBlock error: %s\n", errors.WithStack(err))
			time.Sleep(1 * time.Second)
			continue
		}
		/* Create SyncBlock end */

		ctx.SyncedBlockNumber = block.Number()
		ctx.SyncedBlockHash = common.HexToHash(block.Hash())
	}
}

func rollbackBlock(ctx *svc.ServiceContext) {
	for {
		rollbackBlockNumber := ctx.SyncedBlockNumber

		log.Infof("[Handler.SyncBlock.RollBackBlock]  Try to rollback block number: %d\n", rollbackBlockNumber)

		blockJSON, err := rpc.HTTPPostJSON("", ctx.Config.L2RPCUrl, "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\""+fmt.Sprintf("0x%X", rollbackBlockNumber)+"\", true],\"id\":1}")
		if err != nil {
			log.Errorf("[Handler.SyncBlock.RollRackBlock]Rollback block by number error: %s\n", errors.WithStack(err))
			continue
		}

		rollbackBlock := rpc.ParseJSONBlock(string(blockJSON))
		log.Errorf("[Handler.SyncBlock.RollRackBlock] rollbackBlock: %s, syncedBlockHash: %s \n", rollbackBlock.Hash(), ctx.SyncedBlockHash)

		if common.HexToHash(rollbackBlock.Hash()) == ctx.SyncedBlockHash {
			err = ctx.DB.Transaction(func(tx *gorm.DB) error {
				err = tx.Model(schema.SyncBlock{}).Where(" (status = ? or status = ?) AND block_number>?",
					schema.BlockValid, schema.BlockPending, ctx.SyncedBlockNumber).Update("status", schema.BlockRollback).Error
				if err != nil {
					log.Errorf("[Handler.SyncBlock.RollRackBlock] Rollback Block err: %s\n", errors.WithStack(err))
					return err
				}
				return nil
			})
			if err != nil {
				log.Errorf("[Handler.SyncBlock.RollRackBlock] Rollback db transaction err: %s\n", errors.WithStack(err))
				continue
			}
			log.Infof("[Handler.SyncBlock.RollRackBlock] Rollback blocks is Stop\n")
			return
		}
		var previousBlock schema.SyncBlock
		rest := ctx.DB.Where("block_number = ? AND (status = ? or status = ?) ", rollbackBlockNumber-1, schema.BlockValid, schema.BlockPending).First(&previousBlock)
		if rest.Error != nil {
			log.Errorf("[Handler.RollRackBlock] Previous block by number error: %s\n", errors.WithStack(rest.Error))
			continue
		}
		ctx.SyncedBlockNumber = previousBlock.BlockNumber
		ctx.SyncedBlockHash = common.HexToHash(previousBlock.BlockHash)
	}
}
