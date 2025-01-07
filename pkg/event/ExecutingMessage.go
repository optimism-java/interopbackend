package event

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ExecutingMessageName = "ExecutingMessage"
	ExecutingMessageHash = crypto.Keccak256([]byte("ExecutingMessage(bytes32,(address,uint256,uint256,uint256,uint256))"))
)

type ExecutingMessage struct {
	MsgHash     string   `json:"msgHash"`
	Origin      string   `json:"origin"`
	BlockNumber *big.Int `json:"blockNumber"`
	LogIndex    *big.Int `json:"logIndex"`
	Timestamp   *big.Int `json:"timestamp"`
	ChainId     *big.Int `json:"chainId"`
}

func (*ExecutingMessage) Name() string {
	return ExecutingMessageName
}

func (*ExecutingMessage) EventHash() common.Hash {
	return common.BytesToHash(ExecutingMessageHash)
}

func (t *ExecutingMessage) ToObj(data string) error {
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return err
	}
	return nil
}

func (*ExecutingMessage) Data(log types.Log) (string, error) {
	message := &ExecutingMessage{
		MsgHash:     log.Topics[1].Hex(),                        // indexed msgHash from topics[1]
		Origin:      common.BytesToAddress(log.Data[:32]).Hex(), // first 32 bytes for address
		BlockNumber: new(big.Int).SetBytes(log.Data[32:64]),     // next 32 bytes for blockNumber
		LogIndex:    new(big.Int).SetBytes(log.Data[64:96]),     // next 32 bytes for logIndex
		Timestamp:   new(big.Int).SetBytes(log.Data[96:128]),    // next 32 bytes for timestamp
		ChainId:     new(big.Int).SetBytes(log.Data[128:160]),   // next 32 bytes for chainId
	}

	data, err := ToJSON(message)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (t *ExecutingMessage) GetExecuteMsgHash(log types.Log) string {
	return log.Topics[1].Hex()
}
