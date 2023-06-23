package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/pkg"
	"project/internal/pkg/sqltools"
)

type PostRepository interface {
	GetParentPost(ctx context.Context, post *models.Post) (*models.Post, error)
	UpdatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	GetDetailsPost(ctx context.Context, post *models.Post, params *pkg.PostDetailsParams) (*models.PostDetails, error)
}

type postPostgres struct {
	conn *sql.DB
}

func NewPostPostgres(conn *sql.DB) PostRepository {
	return &postPostgres{
		conn,
	}
}

func (p postPostgres) GetParentPost(ctx context.Context, post *models.Post) (*models.Post, error) {
	res := &models.Post{}

	row := p.conn.QueryRowContext(ctx, `SELECT thread_id
		FROM posts
		WHERE post_id = $1;`, post.Parent)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, pkg.ErrPostParentNotFound
		}

		return nil, row.Err()
	}

	err := row.Scan(&res.Thread)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrPostParentNotFound
		}

		return nil, err
	}

	return res, nil
}

func (p postPostgres) UpdatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	res := &models.Post{}

	err := sqltools.RunTxOnConn(ctx, pkg.TxInsertOptions, p.conn, func(ctx context.Context, tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `UPDATE posts
		SET message   = COALESCE(NULLIF(TRIM($2), ''), message),
			is_edited = CASE
					WHEN TRIM($2) = message THEN is_edited
					ELSE true
				END
		WHERE post_id = $1
		RETURNING parent, author, forum, thread_id, created, message, is_edited;`, post.ID, post.Message)
		if row.Err() != nil {
			return row.Err()
		}

		postTime := time.Time{}

		err := row.Scan(
			&res.Parent,
			&res.Author.Nickname,
			&res.Forum,
			&res.Thread,
			&postTime,
			&res.Message,
			&res.IsEdited)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return pkg.ErrSuchPostNotFound
			}

			return err
		}

		res.Created = postTime.Format(time.RFC3339)

		res.ID = post.ID

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p postPostgres) GetDetailsPost(ctx context.Context, post *models.Post, params *pkg.PostDetailsParams) (*models.PostDetails, error) {
	res := &models.PostDetails{}

	res.Post.ID = post.ID

	row := p.conn.QueryRowContext(ctx, `SELECT parent, author, message, is_edited, forum, thread_id, created
		FROM posts
		WHERE post_id = $1;`, post.ID)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, pkg.ErrSuchPostNotFound
		}

		return nil, row.Err()
	}

	err := row.Scan(
		&res.Post.Parent,
		&res.Post.Author.Nickname,
		&res.Post.Message,
		&res.Post.IsEdited,
		&res.Post.Forum,
		&res.Post.Thread,
		&res.Post.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchPostNotFound
		}

		return nil, err
	}

	for _, value := range params.Related {
		switch value {
		case pkg.PostDetailForum:
			row := p.conn.QueryRowContext(ctx, `SELECT title, users_nickname, slug, posts, threads
				FROM forums 
				WHERE slug = $1;`, res.Post.Forum)
			if row.Err() != nil {
				if errors.Is(row.Err(), sql.ErrNoRows) {
					return nil, pkg.ErrSuchPostNotFound
				}

				return nil, row.Err()
			}

			err := row.Scan(
				&res.Forum.Title,
				&res.Forum.User,
				&res.Forum.Slug,
				&res.Forum.Posts,
				&res.Forum.Threads)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, pkg.ErrSuchPostNotFound
				}

				return nil, err
			}
		case pkg.PostDetailAuthor:
			row := p.conn.QueryRowContext(ctx, `SELECT nickname, fullname, about, email
				FROM users 
				WHERE nickname = $1;`, res.Post.Author.Nickname)
			if row.Err() != nil {
				if errors.Is(row.Err(), sql.ErrNoRows) {
					return nil, pkg.ErrSuchPostNotFound
				}

				return nil, row.Err()
			}

			err := row.Scan(
				&res.Author.Nickname,
				&res.Author.FullName,
				&res.Author.About,
				&res.Author.Email)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, pkg.ErrSuchPostNotFound
				}

				return nil, err
			}
		case pkg.PostDetailThread:
			row := p.conn.QueryRowContext(ctx, `SELECT thread_id, title, author, forum, message, votes, slug, created
				FROM threads
				WHERE thread_id = $1;`, res.Post.Thread)
			if row.Err() != nil {
				if errors.Is(row.Err(), sql.ErrNoRows) {
					return nil, pkg.ErrSuchPostNotFound
				}

				return nil, row.Err()
			}

			err := row.Scan(
				&res.Thread.ID,
				&res.Thread.Title,
				&res.Thread.Author,
				&res.Thread.Forum,
				&res.Thread.Message,
				&res.Thread.Votes,
				&res.Thread.Slug,
				&res.Thread.Created)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, pkg.ErrSuchPostNotFound
				}

				return nil, err
			}
		}
	}

	return res, nil
}
