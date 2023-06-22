package repository

import (
	"context"
	"github.com/jmoiron/sqlx"

	"project/internal/models"
	"project/internal/pkg"
)

type ServiceRepository interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (*models.StatusService, error)
}

type servicePostgres struct {
	db *sqlx.DB
}

func NewServicePostgres(db *sqlx.DB) ServiceRepository {
	return &servicePostgres{
		db,
	}
}

func (s servicePostgres) Clear(ctx context.Context) error {
	row := s.db.QueryRowContext(ctx, `TRUNCATE TABLE forums, posts, threads, user_forums, users, user_votes CASCADE;`)
	if row.Err() != nil {
		return pkg.ErrWorkDatabase
	}

	return nil
}

func (s servicePostgres) GetStatus(ctx context.Context) (*models.StatusService, error) {
	res := &models.StatusService{}

	rowCounters := s.db.QueryRowContext(ctx, `SELECT (SELECT count(*) FROM forums) AS forums,
       (SELECT count(*) FROM posts)  AS posts,
       (SELECT count(*) FROM threads) AS threads,
       (SELECT count(*) FROM users)  AS users`)
	if rowCounters.Err() != nil {
		return nil, pkg.ErrWorkDatabase
	}

	err := rowCounters.Scan(
		&res.Forum,
		&res.Post,
		&res.Thread,
		&res.User)
	if err != nil {
		return nil, err
	}

	return res, nil
}
