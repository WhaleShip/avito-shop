package store

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func CloseAndLogMock(mock pgxmock.PgxConnIface) {
	if err := mock.Close(context.Background()); err != nil {
		log.Println(err)
	}
}

func TxCommitAndLog(tx pgx.Tx) {
	err := tx.Commit(context.Background())
	if err != nil {
		log.Println(err)
	}
}
