package repository

import (
	"context"
	"database/sql"

	"project/internal/models"
	"project/internal/pkg"
	"project/internal/pkg/sqltools"
)

type VoteRepository interface {
	CheckExistVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) (bool, error)
	UpdateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error
	CreateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error
}

type votePostgres struct {
	conn *sql.DB
}

func NewVotePostgres(conn *sql.DB) VoteRepository {
	return &votePostgres{
		conn,
	}
}

func (v votePostgres) CheckExistVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) (bool, error) {
	res := false

	row := v.conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_votes WHERE nickname = $1 AND thread_id = $2);`, params.Nickname, thread.ID)
	if row.Err() != nil {
		return false, row.Err()
	}

	err := row.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (v votePostgres) UpdateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error {
	err := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, v.conn, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `UPDATE user_votes
			SET voice = $3
			WHERE thread_id = $1
			  AND nickname = $2
			  AND voice != $3;`, thread.ID, params.Nickname, params.Voice)
		if row.Err() != nil {
			return row.Err()
		}

		return nil
	})

	return err
}

func (v votePostgres) CreateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error {
	err := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, v.conn, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `INSERT INTO user_votes(nickname, thread_id, voice)
			VALUES ($1, $2, $3);`, params.Nickname, thread.ID, params.Voice)
		if row.Err() != nil {
			return row.Err()
		}

		return nil
	})

	return err
}
