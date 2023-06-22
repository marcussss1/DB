package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/pkg"
)

type ForumRepository interface {
	CheckExistForum(ctx context.Context, forum *models.Forum) (bool, error)
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	GetDetailsForumBySlug(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	GetThreads(ctx context.Context, forum *models.Forum, params *pkg.GetThreadsParams) ([]*models.Thread, error)
	GetUsers(ctx context.Context, forum *models.Forum, params *pkg.GetUsersParams) ([]*models.User, error)
}

type forumPostgres struct {
	db *sqlx.DB
}

func NewForumPostgres(db *sqlx.DB) ForumRepository {
	return &forumPostgres{
		db,
	}
}

func (f forumPostgres) CheckExistForum(ctx context.Context, forum *models.Forum) (bool, error) {
	var res bool

	row := f.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM forums WHERE slug = $1);`, forum.Slug)
	if row.Err() != nil {
		return false, pkg.ErrWorkDatabase
	}

	err := row.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (f forumPostgres) CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	row := f.db.QueryRowContext(ctx, `INSERT INTO forums(title, users_nickname, slug)
		VALUES ($1, $2, $3);`, forum.Title, forum.User, forum.Slug)
	if row.Err() != nil {
		return nil, pkg.ErrWorkDatabase

	}

	return forum, nil
}

func (f forumPostgres) GetDetailsForumBySlug(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	row := f.db.QueryRowContext(ctx, `SELECT title, users_nickname, posts, threads, slug
		FROM forums
		WHERE slug = $1`, forum.Slug)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, pkg.ErrSuchForumNotFound
		}

		return nil, pkg.ErrWorkDatabase

	}

	err := row.Scan(
		&forum.Title,
		&forum.User,
		&forum.Posts,
		&forum.Threads,
		&forum.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchUserNotFound
		}

		return nil, pkg.ErrWorkDatabase
	}

	return forum, nil
}

func (f forumPostgres) GetThreads(ctx context.Context, forum *models.Forum, params *pkg.GetThreadsParams) ([]*models.Thread, error) {
	query := `SELECT t.thread_id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created
		FROM threads AS t
    	LEFT JOIN forums f ON t.forum = f.slug
		WHERE f.slug = $1 `

	orderBy := "ORDER BY t.created "
	querySince := " AND t.created >= $2 "

	var rows *sql.Rows
	var err error

	if params.Desc {
		orderBy += "DESC"
	}

	if params.Limit > 0 {
		orderBy += fmt.Sprintf(" LIMIT %d", params.Limit)
	}

	switch {
	case params.Since != "" && params.Desc:
		querySince = " AND t.created <= $2 "
	case params.Since != "" && !params.Desc:
		querySince = " AND t.created >= $2 "
	}

	var values []interface{}

	if params.Since != "" {
		query += querySince + orderBy

		values = []interface{}{forum.Slug, params.Since}
	} else {
		query += orderBy

		values = []interface{}{forum.Slug}
	}

	res := make([]*models.Thread, 0)

	rows, err = f.db.QueryContext(ctx, query, values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchThreadNotFound
		}

		return nil, pkg.ErrWorkDatabase

	}
	defer rows.Close()

	for rows.Next() {
		thread := &models.Thread{}

		err = rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		if err != nil {
			return nil, err
		}

		res = append(res, thread)
	}

	return res, nil
}

func (f forumPostgres) GetUsers(ctx context.Context, forum *models.Forum, params *pkg.GetUsersParams) ([]*models.User, error) {
	var rows *sql.Rows
	var err error

	query := `SELECT u.nickname, u.fullname, u.about, u.email
		FROM user_forums u
		WHERE u.forum = $1 `

	switch {
	case params.Desc && params.Since != "":
		query += fmt.Sprintf(" AND u.nickname < '%s'", params.Since)
	case params.Since != "":
		query += fmt.Sprintf(" AND u.nickname > '%s'", params.Since)
	}

	query += " ORDER BY u.nickname "

	if params.Desc {
		query += "DESC"
	}

	query += fmt.Sprintf(" LIMIT %d", params.Limit)

	res := make([]*models.User, 0)

	rows, err = f.db.QueryContext(ctx, query, forum.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchThreadNotFound
		}

		return nil, pkg.ErrWorkDatabase
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}

		err = rows.Scan(
			&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email)
		if err != nil {
			return nil, err
		}

		res = append(res, user)
	}

	return res, nil
}
