package route

import (
	"github.com/ZengDaWei/go-ethereum/contract"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventHandleMap 合约事件处理表,一个合约（string key）对应另外一个 map(key => 处理函数)
var EventHandleMap = make(map[string]map[string]func(log types.Log) error)

// ContractEventMap 合约事件表,一个合约（string key）对应多个事件（[]string）
var ContractEventMap = make(map[string][]string)

func Add(contractAddress string, eventSig string, handleFunction func(log types.Log) error) {
	// make 只初始化了第一层 map
	if EventHandleMap[contractAddress] == nil {
		EventHandleMap[contractAddress] = make(map[string]func(log types.Log) error)
	}

	EventHandleMap[contractAddress][eventSig] = handleFunction
	eventSlice := ContractEventMap[contractAddress]
	ContractEventMap[contractAddress] = append(eventSlice, eventSig)

	exists := false

	for _, ars := range contract.ContractAddresses {
		if ars == contractAddress {
			exists = true
		}
	}

	if exists == false {
		contract.ContractAddresses = append(contract.ContractAddresses, contractAddress)
	}
}
