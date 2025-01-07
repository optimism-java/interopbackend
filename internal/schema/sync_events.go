package schema

const (
	EventPending  = "pending"
	EventValid    = "valid"
	EventRollback = "rollback"
	EventInvalid  = "invalid"
)

type SyncEvent struct {
	Base
	SyncBlockID     int64  `json:"sync_block_id"`
	Blockchain      string `json:"blockchain"`
	BlockTime       int64  `json:"block_time"`
	BlockNumber     int64  `json:"block_number"`
	BlockHash       string `json:"block_hash"`
	BlockLogIndexed int64  `json:"block_log_indexed"`
	TxIndex         int64  `json:"tx_index"`
	TxHash          string `json:"tx_hash"`
	EventName       string `json:"event_name"`
	EventHash       string `json:"event_hash"`
	ContractAddress string `json:"contract_address"`
	Data            string `json:"data"`
	Status          string `json:"status"`
	RetryCount      int64  `json:"retry_count"`
	PayloadMsg      string `json:"pay_load_msg"`
	ExecuteMsgHash  string `json:"execute_msg_hash"`
	RelayedMsgHash  string `json:"relayed_msg_hash"`
}

func (SyncEvent) TableName() string {
	return "sync_events"
}
