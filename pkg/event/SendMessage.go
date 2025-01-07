package event

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	SendMessageName = "SentMessage"
	SendMessageHash = crypto.Keccak256([]byte("SentMessage(uint256,address,uint256,address,bytes)"))
)

type SendMessage struct {
	MsgHash     string   `json:"msgHash"`
	Destination *big.Int `json:"destination"`
	Target      string   `json:"target"`
	Nonce       *big.Int `json:"messageNonce"`
	Sender      string   `json:"sender"`
	Message     string   `json:"message"`
}

func (*SendMessage) Name() string {
	return SendMessageName
}

func (*SendMessage) EventHash() common.Hash {
	return common.BytesToHash(SendMessageHash)
}

func (t *SendMessage) ToObj(data string) error {
	err := json.Unmarshal([]byte(data), &t)
	if err != nil {
		return err
	}
	return nil
}

func (*SendMessage) Data(log types.Log) (string, error) {
	message := &SendMessage{
		MsgHash:     log.Topics[0].Hex(),
		Destination: new(big.Int).SetBytes(log.Topics[1].Bytes()),
		Target:      common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		Nonce:       new(big.Int).SetBytes(log.Topics[3].Bytes()),
		Sender:      common.BytesToAddress(log.Data[:32]).Hex(),
		Message:     common.Bytes2Hex(log.Data[32:]),
	}

	data, err := ToJSON(message)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (t *SendMessage) GetExecuteMsgHash(log types.Log) string {
	msg := []byte{}
	for _, topic := range log.Topics {
		msg = append(msg, topic.Bytes()...)
	}
	msg = append(msg, log.Data...)
	s := crypto.Keccak256Hash(msg)
	return s.Hex()
}

func (t *SendMessage) GetRelayedMsgHash(log types.Log, sourceChain int64) (string, error) {
	message := &SendMessage{
		MsgHash:     log.Topics[0].Hex(),
		Destination: new(big.Int).SetBytes(log.Topics[1].Bytes()),
		Target:      common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		Nonce:       new(big.Int).SetBytes(log.Topics[3].Bytes()),
		Sender:      common.BytesToAddress(log.Data[:32]).Hex(),
		Message:     common.Bytes2Hex(log.Data[32:]),
	}

	destination := big.NewInt(message.Destination.Int64())
	source := big.NewInt(sourceChain)
	nonce := big.NewInt(message.Nonce.Int64())
	sender := common.HexToAddress(message.Sender)
	target := common.HexToAddress(message.Target)
	messages := common.FromHex(message.Message)
	msgOffset := new(big.Int).SetBytes(messages[:32])
	msglength := new(big.Int).SetBytes(messages[32:64])
	endMsgOffset := new(big.Int).Add(msgOffset, msglength)
	finalMsg := messages[msgOffset.Int64():endMsgOffset.Int64()]

	arguments := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
		{Type: abi.Type{T: abi.AddressTy}},
		{Type: abi.Type{T: abi.AddressTy}},
		{Type: abi.Type{T: abi.BytesTy}},
	}

	packed, _ := arguments.Pack(destination, source, nonce, sender, target, finalMsg)
	messageHash := crypto.Keccak256Hash(packed)
	return messageHash.Hex(), nil
}
