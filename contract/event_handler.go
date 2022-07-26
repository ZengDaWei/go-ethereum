package contract

import (
	"context"

	"github.com/ZengDaWei/go-ethereum/address"
	"github.com/ZengDaWei/go-ethereum/route"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func HandleEvent(ctx context.Context, log types.Log) error {
	for _, contractAddress := range address.ContractAddresses {
		if log.Address != common.HexToAddress(contractAddress) {
			continue
		}
		for event, handleFunction := range route.EventHandleMap[contractAddress] {
			if crypto.Keccak256Hash([]byte(event)) == log.Topics[0] {
				err := handleFunction(ctx, log)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
