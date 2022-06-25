package contract

import (
	"github.com/ZengDaWei/go-ethereum/route"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func HandleEvent(log types.Log) error {
	for _, contractAddress := range ContractAddresses {
		if log.Address.String() != contractAddress {
			continue
		}
		for event, handleFunction := range route.EventHandleMap[contractAddress] {
			if crypto.Keccak256Hash([]byte(event)) == log.Topics[0] {
				err := handleFunction(log)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
