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
		log.Println(err)
		return nil
	}
	return balance
}

// GetChainId Get chain ID
func GetChainId(client *ethclient.Client) (*big.Int, error) {
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return chainID, nil
}

// GetLatestBlock Get latest block
func GetLatestBlock(client *ethclient.Client) *types.Block {
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return block
}

// GetBlockByNumber Get block
func GetBlockByNumber(client *ethclient.Client, blockNumber *big.Int) *types.Block {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Println(err)
		return nil
	}

	// TODO: Parse block data to human readable json
	return block
}

// GetTransactionByHash Get transaction by hash
func GetTransactionByHash(client *ethclient.Client, txHash string) *types.Transaction {
	hash := common.HexToHash(txHash)
	transaction, _, err := client.TransactionByHash(context.Background(), hash)
	if err != nil {
		log.Println(err)
	}

	// TODO: Parse transaction data to human readable json
	return transaction
}
