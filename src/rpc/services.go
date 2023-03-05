package rpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	own "matsya/src/common"

	"github.com/ethereum/go-ethereum"
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

func parseERC20Logs(log *types.Log) {
	const erc20TransferTopicHash = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	// const erc20TransferTopicHash = "0xa9059cbb2ab09eb219583f4a59a5d0623ade346d962bcd4e46b11da047c9049b"

	if len(log.Topics) > 0 && log.Topics[0].String() == erc20TransferTopicHash {

		if len(log.Topics) == 3 {
			fmt.Println("topic:\t", log.TxIndex, "\t", log.Topics[0])
		}

		for _, topic := range log.Topics {
			fmt.Println("topic:\t", log.TxIndex, "\t", topic.String())
		}
	}

}

func GetBlockLogs(client *ethclient.Client, block *types.Block) *own.Block {
	blockHash := block.Hash()

	blk := &own.Block{
		Hash:      &blockHash,
		Number:    block.Number(),
		Timestamp: block.Time(),
	}

	filterQuery := ethereum.FilterQuery{
		BlockHash: &blockHash,
	}

	blockRecord, transactionRecords := GetBlockTransactions(block)

	fmt.Println(blockRecord)
	fmt.Println(transactionRecords)

	logs, err := client.FilterLogs(context.Background(), filterQuery)
	if err != nil {
		log.Println(err)
	}

	for _, log := range logs {
		fmt.Println("index:\t", log.TxIndex)
		fmt.Println("address:\t", log.Address.String())
		// fmt.Println("data:\t", log.Data)

		parseERC20Logs(&log)
	}

	return blk
}

func GetBlockTransactions(block *types.Block) (*own.Block, []*own.Transaction) {
	blockHash := block.Hash()

	blk := &own.Block{
		Hash:      &blockHash,
		Number:    block.Number(),
		Timestamp: block.Time(),
	}

	transactions := block.Transactions()

	// txns := []*own.Transaction{}
	var txns []*own.Transaction

	fmt.Println("No of transactions ::", len(transactions))

	for _, transaction := range transactions {
		txn := GetTransactionData(block, transaction)
		txns = append(txns, txn)
	}

	return blk, txns
}

func GetTransactionData(block *types.Block, transaction *types.Transaction) *own.Transaction {
	txHash := transaction.Hash()
	fromAddress := GetTransactionMessage(transaction).From().Hash()
	toAddress := transaction.To().Hash()
	txData := hex.EncodeToString(transaction.Data())

	txn := &own.Transaction{
		Hash:           &txHash,
		From:           &fromAddress,
		To:             &toAddress,
		Value:          transaction.Value(),
		GasUsed:        transaction.Gas(),
		GasPrice:       transaction.GasPrice().Uint64(),
		GasFee:         nil,
		Data:           txData,
		BlockNumber:    block.Number(),
		BlockTimestamp: block.Time(),
		ChainID:        transaction.ChainId(),
	}

	return txn
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

	// GetBlockLogs(client, block)

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

// GetTransactionEvents Get transaction events
func GetTransactionEvents(client *ethclient.Client, txHash string) any {
	return nil
}

// GetTransactionLogs Get transaction logs
func GetTransactionLogs(client *ethclient.Client, txHash string) any {
	hash := common.HexToHash(txHash)
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Println(err)
	}

	for _, log := range receipt.Logs {
		fmt.Println("index:\t", log.TxIndex)
		fmt.Println("topics:\t", log.Topics)
		fmt.Println("data:\t", log.Data)
		fmt.Println("topics:\t", log.Address)
	}

	return nil
}

func GetTransactionMessage(tx *types.Transaction) types.Message {
	msg, err := tx.AsMessage(types.LatestSignerForChainID(tx.ChainId()), nil)
	if err != nil {
		log.Fatal(err)
	}
	return msg
}

func ParseTransactionBaseInfo(tx *types.Transaction) map[string]any {
	data := map[string]any{}
	data["Hash"] = tx.Hash().Hex()
	data["ChainId"] = tx.ChainId()
	data["Value"] = tx.Value().String()
	data["From"] = GetTransactionMessage(tx).From().Hex()
	data["To"] = tx.To().Hex()
	data["Gas"] = tx.Gas()
	data["Gas"] = tx.GasPrice().Uint64()
	data["Nonce"] = tx.Nonce()
	data["Transaction"] = hex.EncodeToString(tx.Data())
	return data
}
