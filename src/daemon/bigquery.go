package daemon

import (
	"context"
	"math/big"
	"matsya/src/config"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type BlockData struct {
	Number    int64               `json:"number" bigquery:"number"`
	Timestamp time.Time           `json:"timestamp" bigquery:"timestamp"`
	Hash      bigquery.NullString `json:"hash" bigquery:"hash"`
}

type TransactionData struct {
	Hash                     string              `json:"hash" bigquery:"hash"`
	Nonce                    bigquery.NullInt64  `json:"nonce" bigquery:"nonce"`
	TransactionIndex         bigquery.NullInt64  `json:"transaction_index" bigquery:"transaction_index"`
	FromAddress              bigquery.NullString `json:"from_address" bigquery:"from_address"`
	ToAddress                bigquery.NullString `json:"to_address" bigquery:"to_address"`
	Value                    *big.Rat            `json:"value" bigquery:"value"`
	Gas                      bigquery.NullInt64  `json:"gas" bigquery:"gas"`
	GasPrice                 bigquery.NullInt64  `json:"gas_price" bigquery:"gas_price"`
	Input                    bigquery.NullString `json:"input" bigquery:"input"`
	ReceiptCumulativeGasUsed bigquery.NullInt64  `json:"receipt_cumulative_gas_used" bigquery:"receipt_cumulative_gas_used"`
	ReceiptGasUsed           bigquery.NullInt64  `json:"receipt_gas_used" bigquery:"receipt_gas_used"`
	ReceiptContractAddress   bigquery.NullString `json:"receipt_contract_address,omitempty" bigquery:"receipt_contract_address,nullable"`
	ReceiptRoot              bigquery.NullString `json:"receipt_root" bigquery:"receipt_root"`
	ReceiptStatus            bigquery.NullInt64  `json:"receipt_status" bigquery:"receipt_status"`
	BlockTimestamp           time.Time           `json:"block_timestamp" bigquery:"block_timestamp"`
	BlockNumber              bigquery.NullInt64  `json:"block_number" bigquery:"block_number"`
	BlockHash                bigquery.NullString `json:"block_hash" bigquery:"block_hash"`
	MaxFeePerGas             bigquery.NullInt64  `json:"max_fee_per_gas" bigquery:"max_fee_per_gas"`
	MaxPriorityFeePerGas     bigquery.NullInt64  `json:"max_priority_fee_per_gas" bigquery:"max_priority_fee_per_gas"`
	TransactionType          bigquery.NullInt64  `json:"transaction_type" bigquery:"transaction_type"`
	ReceiptEffectiveGasPrice bigquery.NullInt64  `json:"receipt_effective_gas_price" bigquery:"receipt_effective_gas_price"`
}

func GetBigQueryClient(ctx context.Context) (*bigquery.Client, error) {
	projectID := config.GetString("bigquery.project_id")
	return bigquery.NewClient(ctx, projectID)
}

func GetQueryIterator(
	projectID string,
	query string,
	params []bigquery.QueryParameter,
) (*bigquery.RowIterator, error) {
	// projectID := config.GetString("bigquery.project_id")
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	q := client.Query(query)
	if len(params) > 0 {
		q.Parameters = params
	}

	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func MapBlockData(it *bigquery.RowIterator) ([]*BlockData, error) {
	var blocks []*BlockData
	for {
		var row BlockData
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, &row)
	}
	return blocks, nil
}

func MapTransactionData(it *bigquery.RowIterator) ([]*TransactionData, error) {
	var transactions []*TransactionData
	for {
		var row TransactionData
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &row)
	}
	return transactions, nil
}

func GetTransactionsData(projectID string) ([]*TransactionData, error) {
	query := "SELECT `hash`, `nonce`, `transaction_index`, `from_address`, `to_address`, " +
		"`value`, `gas`, `gas_price`, `input`, `receipt_cumulative_gas_used`, `receipt_gas_used`, `receipt_contract_address`, " +
		"`receipt_root`, `receipt_status`, `block_timestamp`, `block_number`, `block_hash`, `max_fee_per_gas`, " +
		"`max_priority_fee_per_gas`, `transaction_type`, `receipt_effective_gas_price` \n" +
		" FROM `bigquery-public-data.crypto_ethereum.transactions` \n" +
		" WHERE DATE(`block_timestamp`) = ? " +
		"LIMIT 10 "

	var params []bigquery.QueryParameter

	dateString := "2023-03-06"
	param := bigquery.QueryParameter{Value: dateString}

	params = append(params, param)

	it, _ := GetQueryIterator(projectID, query, params)
	return MapTransactionData(it)
}

func GetBlocksData(projectID string) ([]*BlockData, error) {
	query :=
		"SELECT `number`, `hash`, `timestamp` FROM `bigquery-public-data.crypto_ethereum.blocks` " +
			"WHERE DATE(timestamp) = ? " +
			"LIMIT 10 "

	var params []bigquery.QueryParameter

	dateString := "2023-03-06"
	param := bigquery.QueryParameter{Value: dateString}

	params = append(params, param)

	it, _ := GetQueryIterator(projectID, query, params)
	return MapBlockData(it)
}
