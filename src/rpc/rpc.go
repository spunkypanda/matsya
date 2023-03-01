package rpc

import (
	"context"
	"log"
	"math/big"
	"matsya/src/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

func connectNodeProvider() *ethclient.Client {
	rpcHttp := config.GetString("rpc.https")

	client, err := ethclient.Dial(rpcHttp)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func GetWalletBalance(walletAddress string) *big.Int {
	client := connectNodeProvider()
	defer client.Close()

	// Get the balance of an account
	account := common.HexToAddress(walletAddress)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return balance
}

func GetLatestBlock() *types.Block {
	client := connectNodeProvider()
	defer client.Close()

	// Get the latest known block
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return block
}
