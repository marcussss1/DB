package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/pkg"
)

type UserRepository interface {
	CheckFreeEmail(ctx context.Context, user *models.User) (bool, error)
	CreateUser(ctx context.Context, user *models.User) (models.User, error)
	GetUserByEmailOrNickname(ctx context.Context, user *models.User) ([]models.User, error)
	GetUserByNickname(ctx context.Context, user *models.User) (models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (models.User, error)
}

type userPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) UserRepository {
	return &userPostgres{
		db,
	}
}

func (u userPostgres) CheckFreeEmail(ctx context.Context, user *models.User) (bool, error) {
	res := false

	row := u.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`, user.Email)
	if row.Err() != nil {
		return false, pkg.ErrWorkDatabase
	}

	err := row.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (u userPostgres) CreateUser(ctx context.Context, user *models.User) (models.User, error) {
	rowUser := u.db.QueryRowContext(ctx, `INSERT INTO users(nickname, fullname, about, email)
		VALUES ($1, $2, $3, $4);`, user.Nickname, user.FullName, user.About, user.Email)
	if rowUser.Err() != nil {
		return models.User{}, pkg.ErrWorkDatabase
	}

	return *user, nil
}

func (u userPostgres) GetUserByEmailOrNickname(ctx context.Context, user *models.User) ([]models.User, error) {
	res := make([]models.User, 0)

	rowsUsers, err := u.db.QueryContext(ctx, `SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname = $1 OR email = $2;`, user.Nickname, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchUserNotFound
		}

		return nil, pkg.ErrWorkDatabase
	}
	defer rowsUsers.Close()

	for rowsUsers.Next() {
		values := models.User{}

		err = rowsUsers.Scan(
			&values.Nickname,
			&values.FullName,
			&values.About,
			&values.Email)
		if err != nil {
			return nil, err
		}

		res = append(res, values)
	}

	if len(res) == 0 {
		return nil, pkg.ErrSuchUserNotFound
	}

	return res, nil
}

func (u userPostgres) GetUserByNickname(ctx context.Context, user *models.User) (models.User, error) {
	res := models.User{}

	rowUser := u.db.QueryRowContext(ctx, `SELECT fullname, about, email, nickname
		FROM users
		WHERE nickname = $1;`, user.Nickname)
	if rowUser.Err() != nil {
		return models.User{}, pkg.ErrWorkDatabase
	}

	err := rowUser.Scan(
		&res.FullName,
		&res.About,
		&res.Email,
		&res.Nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, pkg.ErrSuchUserNotFound
		}

		return models.User{}, pkg.ErrWorkDatabase
	}

	return res, nil
}

func (u userPostgres) UpdateUser(ctx context.Context, user *models.User) (models.User, error) {
	res := models.User{}

	rowUser := u.db.QueryRowContext(ctx, `UPDATE users
		SET fullname = COALESCE(NULLIF(TRIM($1), ''), fullname),
		about    = COALESCE(NULLIF(TRIM($2), ''), about),
		email    = COALESCE(NULLIF(TRIM($3), ''), email)
		WHERE nickname = $4 RETURNING fullname, about, email, nickname;`, user.FullName, user.About, user.Email, user.Nickname)
	if rowUser.Err() != nil {
		return models.User{}, pkg.ErrUpdateUserDataConflict
	}

	err := rowUser.Scan(
		&res.FullName,
		&res.About,
		&res.Email,
		&res.Nickname)
	if err != nil {
		return models.User{}, err
	}

	return res, nil
}
