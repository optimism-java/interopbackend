package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	config "github.com/optimism-java/interopbackend/internal/types"
	"github.com/optimism-java/interopbackend/pkg/event"
)

type Event interface {
	Name() string
	EventHash() common.Hash
	Data(log types.Log) (string, error)
	ToObj(data string) error
}

var (
	events    = make([]common.Hash, 0)
	contracts = make([]common.Address, 0)
	EventMap  = make(map[common.Hash][]Event, 0)
	EIP1155   = make([]common.Address, 0)
)

func init() {
	Register(&event.SendMessage{})
	Register(&event.RelayedMessage{})
	Register(&event.ExecutingMessage{})
	cfg := config.GetConfig()
	AddContract(cfg.L2toL2CrossDomainMessenger)
	AddContract(cfg.CrossL2Inbox)

}

func Register(event Event) {
	if len(EventMap[event.EventHash()]) == 0 {
		events = append(events, event.EventHash())
	}
	EventMap[event.EventHash()] = append(EventMap[event.EventHash()], event)
}

func AddContract(contract string) {
	contracts = append(contracts, common.HexToAddress(contract))
}

func RemoveContract(contract string) {
	for index, ct := range contracts {
		if ct == common.HexToAddress(contract) {
			contracts = append(contracts[:index], contracts[index+1:]...)
		}
	}
}

func GetContracts() []common.Address {
	return contracts
}

func GetEvents() []common.Hash {
	return events
}

func GetEvent(eventHash common.Hash) Event {
	EventList := EventMap[eventHash]
	Event := EventList[0]
	return Event
}
