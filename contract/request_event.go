package contract

import (
	"context"
	"github.com/ZengDaWei/go-ethereum/address"
	"github.com/ZengDaWei/go-ethereum/route"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	client     *ethclient.Client
	onceClient sync.Once
)

func GetClient(rpcEndpoint string) (*ethclient.Client, error) {
	var err error
	onceClient.Do(func() {
		client, err = ethclient.Dial(rpcEndpoint)
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetEventLogs(from *big.Int, to *big.Int) ([]types.Log, error) {
	var logResult []types.Log

	for _, contractAddress := range address.ContractAddresses {
		var topics []common.Hash
		for _, event := range route.ContractEventMap[contractAddress] {
			topics = append(topics, crypto.Keccak256Hash([]byte(event)))
		}

		logQuery := ethereum.FilterQuery{
			FromBlock: from,
			ToBlock:   to,
			Addresses: []common.Address{common.HexToAddress(contractAddress)},
			Topics:    [][]common.Hash{topics},
		}

		logs, err := client.FilterLogs(context.Background(), logQuery)
		if err != nil {
			return nil, err
		}
		logResult = append(logResult, logs...)
	}

	return logResult, nil
}
