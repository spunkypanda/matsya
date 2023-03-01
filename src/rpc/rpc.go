package rpc

import (
	"log"
	"matsya/src/config"

	"github.com/ethereum/go-ethereum/ethclient"

	_ "github.com/lib/pq"
)

func GetNodeProvider() *ethclient.Client {
	rpcHttp := config.GetString("rpc.https")

	client, err := ethclient.Dial(rpcHttp)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
