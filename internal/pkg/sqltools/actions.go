package sqltools

import (
	"context"
	"database/sql"
)

func RunTxOnConn(ctx context.Context, options *sql.TxOptions, db *sql.DB, action func(ctx context.Context, tx *sql.Tx) error) error {
	conn, _ := db.Conn(ctx)
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, options)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}()

	err = action(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
