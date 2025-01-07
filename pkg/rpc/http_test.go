package rpc

import (
	"fmt"
	"testing"
)

func TestEthGetBlockByNumber(t *testing.T) {
	json, err := HTTPPostJSON("", "http://35.165.132.220:8545", "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"48203\", true],\"id\":1}")
	if err != nil {
		return
	}
	block := ParseJSONBlock(string(json))
	fmt.Println(block.Timestamp())
	fmt.Println(block.Number())
	fmt.Println(block.Hash())
}
