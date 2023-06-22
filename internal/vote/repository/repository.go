package repository

import (
	"context"
	"github.com/jmoiron/sqlx"

	"project/internal/models"
	"project/internal/pkg"
)

type VoteRepository interface {
	CheckExistVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) (bool, error)
	UpdateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error
	CreateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error
}

type votePostgres struct {
	db *sqlx.DB
}

func NewVotePostgres(db *sqlx.DB) VoteRepository {
	return &votePostgres{
		db,
	}
}

func (v votePostgres) CheckExistVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) (bool, error) {
	res := false

	rowExist := v.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_votes WHERE nickname = $1 AND thread_id = $2);`, params.Nickname, thread.ID)
	if rowExist.Err() != nil {
		return false, pkg.ErrWorkDatabase
	}

	err := rowExist.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (v votePostgres) UpdateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error {
	rowUpdate := v.db.QueryRowContext(ctx, `UPDATE user_votes
		SET voice = $3
		WHERE thread_id = $1 AND nickname = $2 AND voice != $3;`, thread.ID, params.Nickname, params.Voice)
	if rowUpdate.Err() != nil {
		return pkg.ErrWorkDatabase
	}

	return nil
}

func (v votePostgres) CreateVote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) error {
	rowCreate := v.db.QueryRowContext(ctx, `INSERT INTO user_votes(nickname, thread_id, voice)
		VALUES ($1, $2, $3);`, params.Nickname, thread.ID, params.Voice)
	if rowCreate.Err() != nil {
		return pkg.ErrWorkDatabase
	}

	return nil
}
