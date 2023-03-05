package common

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type CustomHandler func(http.ResponseWriter, *http.Request, context.Context)

type Chain struct {
	ID                 int            `db:"id"`
	Name               string         `db:"name"`
	Icon               string         `db:"icon"`
	Currency           string         `db:"currency"`
	BinanceNetworkId   string         `db:"binanceNetworkId"`
	RpcHttp            string         `db:"rpcHttp"`
	IsEIP1559Supported bool           `db:"isEIP1559Supported"`
	Mainnet            bool           `db:"mainnet"`
	Sequence           int            `db:"sequence"`
	CreatedAt          *time.Time     `db:"createdAt"`
	UpdatedAt          *time.Time     `db:"updatedAt"`
	NativeAssetId      sql.NullString `db:"nativeAssetId"`
	DeletedAt          *time.Time     `db:"deletedAt"`
	GnosisSafeApiUrl   sql.NullString `db:"gnosisSafeApiUrl"`
	RpcWs              sql.NullString `db:"rpcWs"`
}

// Contract types
const (
	ERC20   = "erc20"
	ERC721  = "erc721"
	ERC1155 = "erc1155"
)

// Contract
type Contract struct {
	Address *common.Hash
	Kind    string
}

// Currency
type Currency struct {
	Contract    *common.Hash
	Name        string
	Symbol      string
	Decimals    int8
	Metadata    json.RawMessage
	ChainID     int8
	CoingeckoID int8
}

// Block
type Block struct {
	Hash      *common.Hash
	Number    *big.Int
	Timestamp uint64
}

// Transaction record
type Transaction struct {
	Hash           *common.Hash
	From           *common.Hash
	To             *common.Hash
	Value          *big.Int
	ChainID        *big.Int
	BlockNumber    *big.Int
	BlockTimestamp uint64
	GasUsed        uint64
	GasPrice       uint64
	GasFee         *big.Int
	Data           string
}
