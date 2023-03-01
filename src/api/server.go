package api

import (
	"context"
	"encoding/json"
	"fmt"
	"matsya/src/common"
	"matsya/src/rpc"
	"matsya/src/services"
	"net/http"
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

	balance := rpc.GetWalletBalance(walletAddress)
	value := rpc.ToDecimal(balance, 18)

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{"balance": value, "currency": "eth"}
	json.NewEncoder(w).Encode(response)
}

func ChainsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	db := common.GetConnection()
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

	fmt.Println("Listening on ", port)
	http.ListenAndServe(port, nil)
}
