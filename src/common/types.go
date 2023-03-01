package common

import (
	"context"
	"database/sql"
	"net/http"
	"time"
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
