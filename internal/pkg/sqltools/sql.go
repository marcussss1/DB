package sqltools

import (
	"context"
	"database/sql"
	// justifying it
	_ "github.com/jackc/pgx/stdlib"
)

func InsertBatch(ctx context.Context, db *sql.DB, query string, values []interface{}) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
