package handler

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/optimism-java/interopbankend/internal/blockchain"
	"github.com/optimism-java/interopbankend/internal/schema"
	"github.com/optimism-java/interopbankend/internal/svc"
	"github.com/optimism-java/interopbankend/pkg/event"
	"github.com/optimism-java/interopbankend/pkg/log"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

func LogBatchFilter(ctx *svc.ServiceContext, startBlock, endBlock int64, addresses []common.Address, topics [][]common.Hash) ([]*schema.SyncEvent, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(startBlock),
		ToBlock:   big.NewInt(endBlock),
		Topics:    topics,
		Addresses: addresses,
	}
	logs, err := ctx.L2RPC.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return LogsToEvents(ctx, logs, startBlock)
}

func LogFilter(ctx *svc.ServiceContext, block schema.SyncBlock, addresses []common.Address, topics [][]common.Hash) ([]*schema.SyncEvent, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(block.BlockNumber),
		ToBlock:   big.NewInt(block.BlockNumber),
		Topics:    topics,
		Addresses: addresses,
	}
	logs, err := ctx.L2RPC.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Infof("[CancelOrder.Handle] Cancel Pending List Length is %d ,block number is %d \n", len(logs), block.BlockNumber)
	return LogsToEvents(ctx, logs, block.ID)
}

func LogsToEvents(ctx *svc.ServiceContext, logs []types.Log, syncBlockId int64) ([]*schema.SyncEvent, error) {
	var events []*schema.SyncEvent
	blockTimes := make(map[int64]int64)
	sendMessageEvent := &event.SendMessage{}
	executingMessageEvent := &event.ExecutingMessage{}
	relayedMessageEvent := &event.RelayedMessage{}
	for _, vlog := range logs {
		eventHash := event.TopicToHash(vlog, 0)
		contractAddress := vlog.Address
		Event := blockchain.GetEvent(eventHash)
		if Event == nil {
			log.Infof("[LogsToEvents] logs[txHash: %s, contractAddress:%s, eventHash: %s]\n", vlog.TxHash, strings.ToLower(contractAddress.Hex()), eventHash)
			continue
		}

		blockTime := blockTimes[int64(vlog.BlockNumber)]
		if blockTime == 0 {
			block, err := ctx.L2RPC.BlockByNumber(context.Background(), big.NewInt(int64(vlog.BlockNumber)))
			if err != nil {
				return nil, errors.WithStack(err)
			}
			blockTime = int64(block.Time())
		}
		data, err := Event.Data(vlog)
		if err != nil {
			log.Errorf("[LogsToEvents] logs[txHash: %s, contractAddress:%s, eventHash: %s]\n", vlog.TxHash, strings.ToLower(contractAddress.Hex()), eventHash)
			log.Errorf("[LogsToEvents] data err: %s\n", errors.WithStack(err))
			continue
		}
		payloadMsgBytes := payloadBytes(&vlog)
		evt := &schema.SyncEvent{
			Blockchain:      ctx.Config.Blockchain,
			SyncBlockID:     syncBlockId,
			BlockTime:       blockTime,
			BlockNumber:     int64(vlog.BlockNumber),
			BlockHash:       vlog.BlockHash.Hex(),
			BlockLogIndexed: int64(vlog.Index),
			TxIndex:         int64(vlog.TxIndex),
			TxHash:          vlog.TxHash.Hex(),
			EventName:       Event.Name(),
			EventHash:       eventHash.Hex(),
			ContractAddress: strings.ToLower(contractAddress.Hex()),
			Data:            data,
			Status:          schema.EventPending,
			PayloadMsg:      common.Bytes2Hex(payloadMsgBytes),
		}
		if sendMessageEvent.Name() == evt.EventName && sendMessageEvent.EventHash().Hex() == evt.EventHash {
			evt.ExecuteMsgHash = sendMessageEvent.GetExecuteMsgHash(vlog)
			hash, err := sendMessageEvent.GetRelayedMsgHash(vlog, ctx.Config.BlockChainID)
			if err != nil {
				return nil, err
			}
			evt.RelayedMsgHash = hash
		}
		if executingMessageEvent.Name() == evt.EventName && executingMessageEvent.EventHash().Hex() == evt.EventHash {
			evt.ExecuteMsgHash = executingMessageEvent.GetExecuteMsgHash(vlog)
		}
		if relayedMessageEvent.Name() == evt.EventName && relayedMessageEvent.EventHash().Hex() == evt.EventHash {
			evt.RelayedMsgHash = relayedMessageEvent.GetRelayedMessage(vlog)
		}

		events = append(events, evt)
	}
	return events, nil
}

func payloadBytes(log *types.Log) []byte {
	msg := []byte{}
	for _, topic := range log.Topics {
		msg = append(msg, topic.Bytes()...)
	}
	return append(msg, log.Data...)
}
