package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"matsya/src/common"
	"matsya/src/rpc"
	"matsya/src/services"
	"net/http"
	"strconv"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	body := fmt.Sprintln("Hello world")
	fmt.Fprint(w, body)
}

func WalletBalanceHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	walletAddress := r.URL.Query().Get("wallet")

	if walletAddress == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]interface{}{"message": "Missing wallet address in query"}
		json.NewEncoder(w).Encode(response)
		return
	}

	client := rpc.GetNodeProvider()
	defer client.Close()

	balance := rpc.GetWalletBalance(client, walletAddress)
	value := rpc.ToDecimal(balance, 18)

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{"balance": value, "currency": "eth"}
	json.NewEncoder(w).Encode(response)
}

func BlockHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	blockNoString := r.URL.Query().Get("block")

	if blockNoString == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]interface{}{"message": "Missing block in query"}
		json.NewEncoder(w).Encode(response)
		return
	}

	client := rpc.GetNodeProvider()
	defer client.Close()

	blockNo, err := strconv.Atoi(blockNoString)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]interface{}{"message": "Invalid input for block"}
		json.NewEncoder(w).Encode(response)
		return
	}

	block := big.NewInt(int64(blockNo))

	blockData := rpc.GetBlockByNumber(client, block)

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{"block": blockData}
	json.NewEncoder(w).Encode(response)
}

func TransactionHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	hash := r.URL.Query().Get("hash")

	if hash == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]interface{}{"message": "Missing hash in query"}
		json.NewEncoder(w).Encode(response)
		return
	}

	client := rpc.GetNodeProvider()
	defer client.Close()

	transaction := rpc.GetTransactionByHash(client, hash)

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{"transaction": transaction}
	json.NewEncoder(w).Encode(response)
}

func ChainsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	db, err := common.GetConnection()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := map[string]interface{}{"message": "Missing wallet address in query"}
		json.NewEncoder(w).Encode(response)
		return
	}

	chains := services.GetChainsByEnv(db, true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chains)
}

func withContext(handler common.CustomHandler, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx = context.WithValue(ctx, "query", r.URL.Query())
		handler(w, r, ctx)
	}
}

func Serve(port string) {
	ctx := context.Background()

	http.HandleFunc("/", withContext(HomeHandler, ctx))
	http.HandleFunc("/chains", withContext(ChainsHandler, ctx))
	http.HandleFunc("/balance", withContext(WalletBalanceHandler, ctx))
	http.HandleFunc("/transaction", withContext(TransactionHandler, ctx))
	http.HandleFunc("/block", withContext(BlockHandler, ctx))

	fmt.Println("Listening on ", port)
	http.ListenAndServe(port, nil)
}
