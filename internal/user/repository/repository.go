package repository

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/pkg"
	"project/internal/pkg/sqltools"
)

type UserRepository interface {
	CheckFreeEmail(ctx context.Context, user *models.User) (bool, error)
	CreateUser(ctx context.Context, user *models.User) (models.User, error)
	GetUserByEmailOrNickname(ctx context.Context, user *models.User) ([]models.User, error)
	GetUserByNickname(ctx context.Context, user *models.User) (models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (models.User, error)
}

type userPostgres struct {
	conn *sql.DB
}

func NewUserPostgres(conn *sql.DB) UserRepository {
	return &userPostgres{
		conn,
	}
}

func (u userPostgres) CheckFreeEmail(ctx context.Context, user *models.User) (bool, error) {
	res := false

	row := u.conn.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`, user.Email)
	if row.Err() != nil {
		return false, row.Err()
	}

	err := row.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (u userPostgres) CreateUser(ctx context.Context, user *models.User) (models.User, error) {
	err := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, u.conn, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `INSERT INTO users(nickname, fullname, about, email)
			VALUES ($1, $2, $3, $4);`, user.Nickname, user.FullName, user.About, user.Email)
		if row.Err() != nil {
			return row.Err()
		}

		return nil
	})
	if err != nil {
		return models.User{}, err
	}

	return *user, nil
}

func (u userPostgres) GetUserByEmailOrNickname(ctx context.Context, user *models.User) ([]models.User, error) {
	res := make([]models.User, 0)

	row, err := u.conn.QueryContext(ctx, `SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname = $1
		   OR email = $2;`, user.Nickname, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchUserNotFound
		}

		return nil, err
	}
	defer row.Close()

	for row.Next() {
		values := models.User{}

		err = row.Scan(
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

	row := u.conn.QueryRowContext(ctx, `SELECT fullname, about, email, nickname
		FROM users
		WHERE nickname = $1;`, user.Nickname)
	if row.Err() != nil {
		return models.User{}, row.Err()
	}

	err := row.Scan(
		&res.FullName,
		&res.About,
		&res.Email,
		&res.Nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, pkg.ErrSuchUserNotFound
		}

		return models.User{}, err
	}

	return res, nil
}

func (u userPostgres) UpdateUser(ctx context.Context, user *models.User) (models.User, error) {
	res := models.User{}

	err := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, u.conn, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `UPDATE users
			SET fullname = COALESCE(NULLIF(TRIM($1), ''), fullname),
				about    = COALESCE(NULLIF(TRIM($2), ''), about),
				email    = COALESCE(NULLIF(TRIM($3), ''), email)
			WHERE nickname = $4 RETURNING fullname, about, email, nickname;`, user.FullName, user.About, user.Email, user.Nickname)
		if row.Err() != nil {
			return row.Err()
		}

		err := row.Scan(
			&res.FullName,
			&res.About,
			&res.Email,
			&res.Nickname)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return models.User{}, err
	}

	return res, nil
}
