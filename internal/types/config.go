package types

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	// debug", "info", "warn", "error", "panic", "fatal"
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	// "console","json"
	LogFormat                  string `env:"LOG_FORMAT" envDefault:"console"`
	MySQLDataSource            string `env:"MYSQL_DATA_SOURCE" envDefault:"root:root@tcp(127.0.0.1:3367)/OPChainB?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"`
	MySQLMaxIdleConns          int    `env:"MYSQL_MAX_IDLE_CONNS" envDefault:"10"`
	MySQLMaxOpenConns          int    `env:"MYSQL_MAX_OPEN_CONNS" envDefault:"20"`
	MySQLConnMaxLifetime       int    `env:"MYSQL_CONN_MAX_LIFETIME" envDefault:"3600"`
	Blockchain                 string `env:"BLOCKCHAIN" envDefault:"OPChainB"`
	BlockChainID               int64  `env:"BLOCKCHAIN_ID" envDefault:"902"`
	L2RPCUrl                   string `env:"L2_RPC_URL" envDefault:"http://127.0.0.1:9546"`
	RPCRateLimit               int    `env:"RPC_RATE_LIMIT" envDefault:"15"`
	RPCRateBurst               int    `env:"RPC_RATE_BURST" envDefault:"5"`
	FromBlockNumber            int64  `env:"FROM_BLOCK_NUMBER" envDefault:"-1"`
	FromBlockHash              string `env:"FROM_BLOCK_HASH" envDefault:"0x0000000000000000000000000000000000000000000000000000000000000000"`
	L2toL2CrossDomainMessenger string `env:"L2_TO_L2_CROSS_DOMAIN_MESSENGER" envDefault:"0x4200000000000000000000000000000000000023"`
	CrossL2Inbox               string `env:"CROSS_L2_INBOX" envDefault:"0x4200000000000000000000000000000000000022"`
	APIPort                    string `env:"API_PORT" envDefault:"8088"`
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		cfg := &Config{}
		if err := env.Parse(cfg); err != nil {
			log.Panicf("parse config err: %s\n", err)
			return nil
		}
		config = cfg
	}
	return config
}
