package services

import (
	"fmt"

	common "matsya/src/common"

	"github.com/jmoiron/sqlx"
)

func GetChains(db *sqlx.DB) []common.Chain {
	chains := []common.Chain{}

	query := `SELECT * FROM "chain"`

	err := db.Select(&chains, query)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return chains
}

func GetChainsByEnv(db *sqlx.DB, isMainnet bool) []common.Chain {
	chains := []common.Chain{}

	args := common.Chain{Mainnet: isMainnet}

	namedStmt, err2 := db.PrepareNamed(`SELECT * FROM chain WHERE mainnet = :mainnet`)
	if err2 != nil {
		fmt.Println(err2.Error())
		return nil
	}

	err2 = namedStmt.Select(&chains, args)
	if err2 != nil {
		fmt.Println(err2.Error())
		return nil
	}

	return chains
}
