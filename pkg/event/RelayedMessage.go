package event

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	RelayedMessageName = "RelayedMessage"

	RelayedMessageHash = crypto.Keccak256([]byte("RelayedMessage(uint256,uint256,bytes32)"))
)

type RelayedMessage struct {
	Source       *big.Int `json:"source"`
	MessageNonce *big.Int `json:"messageNonce"`
	MessageHash  string   `json:"messageHash"`
}

func (*RelayedMessage) Name() string {
	return RelayedMessageName
}

func (*RelayedMessage) EventHash() common.Hash {
	return common.BytesToHash(RelayedMessageHash)
}

func (t *RelayedMessage) ToObj(data string) error {
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return err
	}
	return nil
}

func (*RelayedMessage) Data(log types.Log) (string, error) {
	transfer := &RelayedMessage{
		Source:       big.NewInt(TopicToInt64(log, 1)),
		MessageNonce: big.NewInt(TopicToInt64(log, 2)),
		MessageHash:  TopicToHash(log, 3).String(),
	}
	data, err := ToJSON(transfer)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (t *RelayedMessage) GetRelayedMessage(log types.Log) string {
	return TopicToHash(log, 3).String()
}
