package rpc

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetWalletBalance Get the balance of an account
func GetWalletBalance(client *ethclient.Client, walletAddress string) *big.Int {
	account := common.HexToAddress(walletAddress)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return balance
}

// GetLatestBlock Get latest block
func GetLatestBlock(client *ethclient.Client) *types.Block {
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return block
}
