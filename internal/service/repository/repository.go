package repository

import (
	"context"
	"database/sql"

	"project/internal/models"
	"project/internal/pkg"
	"project/internal/pkg/sqltools"
)

type ServiceRepository interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (*models.StatusService, error)
}

type servicePostgres struct {
	conn *sql.DB
}

func NewServicePostgres(conn *sql.DB) ServiceRepository {
	return &servicePostgres{
		conn,
	}
}

func (s servicePostgres) Clear(ctx context.Context) error {
	err := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, s.conn, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `TRUNCATE TABLE forums, posts, threads, user_forums, users, user_votes CASCADE;`)
		if row.Err() != nil {
			return row.Err()
		}

		return nil
	})

	return err
}

func (s servicePostgres) GetStatus(ctx context.Context) (*models.StatusService, error) {
	res := &models.StatusService{}

	row := s.conn.QueryRowContext(ctx, `SELECT (SELECT count(*) FROM forums) AS forums,
       (SELECT count(*) FROM posts)  AS posts,
       (SELECT count(*) FROM threads) AS threads,
       (SELECT count(*) FROM users)  AS users`)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err := row.Scan(
		&res.Forum,
		&res.Post,
		&res.Thread,
		&res.User)
	if err != nil {
		return nil, err
	}

	return res, nil
}
