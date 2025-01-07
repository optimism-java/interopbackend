package v0

import migration "github.com/optimism-java/interopbackend/migration/version"

type SyncBlock struct {
	migration.Base
	Blockchain   string `json:"blockchain" gorm:"type:varchar(32);notnull"`
	BlockchainID int64  `json:"blockchain_id" gorm:"type:bigint;notnull"`
	Miner        string `json:"miner" gorm:"type:varchar(42);notnull"`
	BlockTime    int64  `json:"block_time" gorm:"type:bigint;notnull"`
	BlockNumber  int64  `json:"block_number" gorm:"type:bigint;notnull"`
	BlockHash    string `json:"block_hash" gorm:"type:varchar(66);notnull"`
	TxCount      int64  `json:"tx_count" gorm:"type:bigint;notnull;index:tx_count"`
	EventCount   int64  `json:"event_count" gorm:"type:bigint;notnull"`
	ParentHash   string `json:"parent_hash" gorm:"type:varchar(66);notnull"`
	Status       string `json:"status" gorm:"type:varchar(32);notnull;index:status_index"`
	CheckCount   int64  `json:"check_count" gorm:"type:bigint;notnull;index:check_count"`
	HasCrossTx   bool   `json:"has_cross_tx" gorm:"type:tinyint(1);notnull"`
	BlockState   string `json:"block_state" gorm:"type:varchar(32);notnull;index:block_state_index"`
}

func (SyncBlock) TableName() string {
	return "sync_blocks"
}

type SyncEvent struct {
	migration.Base
	SyncBlockID     int64  `json:"sync_block_id" gorm:"type:bigint;notnull"`
	Blockchain      string `json:"blockchain" gorm:"type:varchar(32);notnull"`
	BlockTime       int64  `json:"block_time" gorm:"type:bigint;notnull"`
	BlockNumber     int64  `json:"block_number" gorm:"type:bigint;notnull"`
	BlockHash       string `json:"block_hash" gorm:"type:varchar(66);notnull"`
	BlockLogIndexed int64  `json:"block_log_indexed" gorm:"type:bigint;notnull"`
	TxIndex         int64  `json:"tx_index" gorm:"type:bigint;notnull"`
	TxHash          string `json:"tx_hash" gorm:"type:varchar(66);notnull"`
	EventName       string `json:"event_name" gorm:"type:varchar(32);notnull"`
	EventHash       string `json:"event_hash" gorm:"type:varchar(66);notnull"`
	ContractAddress string `json:"contract_address" gorm:"type:varchar(42);notnull"`
	Data            string `json:"data" gorm:"type:mediumtext;notnull"`
	Status          string `json:"status" gorm:"type:varchar(32);notnull;index:status_index"`
	RetryCount      int64  `json:"retry_count" gorm:"type:bigint;notnull"`
	PayloadMsg      string `json:"pay_load_msg" gorm:"type:mediumtext;notnull"`
	ExecuteMsgHash  string `json:"execute_msg_hash" gorm:"type:varchar(66)"`
	RelayedMsgHash  string `json:"relayed_msg_hash" gorm:"type:varchar(66)"`
}

func (SyncEvent) TableName() string {
	return "sync_events"
}

var ModelSchemaList = []interface{}{SyncBlock{}, SyncEvent{}}
