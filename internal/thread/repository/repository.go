package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/pkg"
	"project/internal/pkg/sqltools"
)

type ThreadRepository interface {
	CreateThread(ctx context.Context, thread *models.Thread) (models.Thread, error)
	CreatePostsByID(ctx context.Context, thread *models.Thread, posts []*models.Post) ([]models.Post, error)
	GetDetailsThreadByID(ctx context.Context, thread *models.Thread) (models.Thread, error)
	GetDetailsThreadBySlug(ctx context.Context, thread *models.Thread) (models.Thread, error)
	UpdateThreadByID(ctx context.Context, thread *models.Thread) (models.Thread, error)

	GetPostsByIDFlat(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error)
	GetPostsByIDTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error)
	GetPostsByIDParentTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error)
}

type threadPostgres struct {
	db *sqlx.DB
}

func NewThreadPostgres(db *sqlx.DB) ThreadRepository {
	return &threadPostgres{
		db,
	}
}

func (t threadPostgres) CreateThread(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	if thread.Created == "" {
		thread.Created = time.Now().Format(time.RFC3339)
	}

	rowThread := t.db.QueryRowContext(ctx, `INSERT INTO threads(title, author, forum, message, slug, created)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING thread_id;`, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created)
	if rowThread.Err() != nil {
		return models.Thread{}, rowThread.Err()
	}

	err := rowThread.Scan(&thread.ID)
	if err != nil {
		return models.Thread{}, err
	}

	if thread.Forum == thread.Slug {
		thread.Slug = ""
	}

	return *thread, nil
}

func (t threadPostgres) CreatePostsByID(ctx context.Context, thread *models.Thread, posts []*models.Post) ([]models.Post, error) {
	query := `INSERT INTO posts(parent, author, message, forum, thread_id, created) VALUES`

	countAttributes := strings.Count(query, ",") + 1

	pos := 0

	countInserts := len(posts)

	values := make([]interface{}, countInserts*countAttributes)

	insertTimeString := time.Now().Format(time.RFC3339)

	for i := 0; i < len(posts); i++ {
		values[pos] = posts[i].Parent
		pos++
		values[pos] = posts[i].Author.Nickname
		pos++
		values[pos] = posts[i].Message
		pos++
		values[pos] = thread.Forum
		pos++
		values[pos] = thread.ID
		pos++
		values[pos] = insertTimeString
		pos++
	}

	insertStatement := sqltools.CreateFullQuery(query, countInserts, countAttributes)

	insertStatement += " RETURNING post_id;"

	rows, err := t.db.QueryContext(ctx, insertStatement, values...)
	if err != nil {
		return nil, err
	}

	res := make([]models.Post, len(posts))

	i := 0
	for rows.Next() {
		err = rows.Scan(&res[i].ID)
		if err != nil {
			return nil, err
		}

		res[i].Created = insertTimeString
		res[i].Parent = posts[i].Parent
		res[i].Author.Nickname = posts[i].Author.Nickname
		res[i].Message = posts[i].Message
		res[i].Forum = thread.Forum
		res[i].Thread = thread.ID

		i++
	}

	return res, nil
}

func (t threadPostgres) GetDetailsThreadByID(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	res := models.Thread{}

	row := t.db.QueryRowContext(ctx, `SELECT title, author, forum, message, votes, slug, created
		FROM threads WHERE thread_id = $1;`, thread.ID)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return models.Thread{}, pkg.ErrSuchThreadNotFound
		}

		return models.Thread{}, row.Err()
	}

	err := row.Scan(
		&res.Title,
		&res.Author,
		&res.Forum,
		&res.Message,
		&res.Votes,
		&res.Slug,
		&res.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Thread{}, pkg.ErrSuchThreadNotFound
		}

		return models.Thread{}, err
	}

	res.ID = thread.ID

	return res, nil
}

func (t threadPostgres) GetDetailsThreadBySlug(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	res := models.Thread{}

	row := t.db.QueryRowContext(ctx, `SELECT thread_id, title, author, forum, message, votes, slug, created
		FROM threads WHERE slug = $1;`, thread.Slug)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return models.Thread{}, pkg.ErrSuchThreadNotFound
		}

		return models.Thread{}, row.Err()
	}

	err := row.Scan(
		&res.ID,
		&res.Title,
		&res.Author,
		&res.Forum,
		&res.Message,
		&res.Votes,
		&res.Slug,
		&res.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Thread{}, pkg.ErrSuchThreadNotFound
		}

		return models.Thread{}, err
	}

	return res, nil
}

func (t threadPostgres) UpdateThreadByID(ctx context.Context, thread *models.Thread) (models.Thread, error) {
	res := models.Thread{}

	rowThread := t.db.QueryRowContext(ctx, `UPDATE threads
		SET title   = COALESCE(NULLIF(TRIM($2), ''), title),
    	message = COALESCE(NULLIF(TRIM($3), ''), message)
		WHERE thread_id = $1
		RETURNING author, forum, votes, slug, created, title, message;`, thread.ID, thread.Title, thread.Message)
	if rowThread.Err() != nil {
		return models.Thread{}, rowThread.Err()
	}

	err := rowThread.Scan(
		&res.Author,
		&res.Forum,
		&res.Votes,
		&res.Slug,
		&res.Created,
		&res.Title,
		&res.Message)
	if err != nil {
		return models.Thread{}, err
	}

	res.ID = thread.ID

	return res, nil
}

func (t threadPostgres) GetPostsByIDFlat(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error) {
	var rows *sql.Rows
	var err error

	query := `SELECT post_id, parent, author, message, is_edited, forum, created FROM posts WHERE thread_id = $1 `

	var values []interface{}

	switch {
	case params.Since != -1 && params.Desc:
		query += " AND post_id < $2"
	case params.Since != -1 && !params.Desc:
		query += " AND post_id > $2"
	case params.Since != -1:
		query += " AND post_id > $2"
	}

	switch {
	case params.Desc:
		query += " ORDER BY created DESC, post_id DESC"
	case !params.Desc:
		query += " ORDER BY created ASC, post_id"
	default:
		query += " ORDER BY created, post_id"
	}

	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", params.Limit)

	if params.Since == -1 {
		values = []interface{}{thread.ID}
	} else {
		values = []interface{}{thread.ID, params.Since}
	}

	res := make([]models.Post, 0)

	rows, err = t.db.QueryContext(ctx, query, values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchPostNotFound
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		timeTmp := time.Time{}

		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author.Nickname,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&timeTmp)
		if err != nil {
			return nil, err
		}

		post.Thread = thread.ID

		post.Created = timeTmp.Format(time.RFC3339)

		res = append(res, post)
	}

	return res, nil
}

func (t threadPostgres) GetPostsByIDTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error) {
	var rows *sql.Rows
	var err error

	query := `SELECT post_id, parent, author, message, is_edited, forum, created FROM posts WHERE thread_id = $1 `

	switch {
	case params.Since != -1 && params.Desc:
		query += " AND path < "
	case params.Since != -1 && !params.Desc:
		query += " AND path > "
	case params.Since != -1:
		query += " AND path > "
	}

	if params.Since != -1 {
		query += fmt.Sprintf(` (SELECT path FROM posts WHERE post_id = %d) `, params.Since)
	}

	switch {
	case params.Desc:
		query += " ORDER BY path DESC"
	case !params.Desc:
		query += " ORDER BY path ASC, post_id"
	default:
		query += " ORDER BY path, post_id"
	}

	query += fmt.Sprintf(" LIMIT NULLIF(%d, 0)", params.Limit)

	res := make([]models.Post, 0)

	rows, err = t.db.QueryContext(ctx, query, thread.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchPostNotFound
		}

		return nil, pkg.ErrWorkDatabase
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		timeTmp := time.Time{}

		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author.Nickname,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&timeTmp)
		if err != nil {
			return nil, err
		}

		post.Thread = thread.ID

		post.Created = timeTmp.Format(time.RFC3339)

		res = append(res, post)
	}

	return res, nil
}

func (t threadPostgres) GetPostsByIDParentTree(ctx context.Context, thread *models.Thread, params *pkg.GetPostsParams) ([]models.Post, error) {
	var rows *sql.Rows
	var err error

	query := ""

	var values []interface{}

	if params.Since == -1 {
		if params.Desc {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 ORDER BY post_id DESC LIMIT $2)
					ORDER BY path[1] DESC, path ASC, post_id ASC;`
		} else {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 ORDER BY post_id ASC LIMIT $2)
					ORDER BY path ASC, post_id ASC;`
		}

		values = []interface{}{thread.ID, params.Limit}
	} else {
		if params.Desc {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 AND path[1] <
					(SELECT path[1] FROM posts WHERE post_id = $2) ORDER BY post_id DESC LIMIT $3)
					ORDER BY path[1] DESC, path ASC, post_id ASC;`
		} else {
			query = `
					SELECT post_id, parent, author, message, is_edited, forum, created FROM posts
					WHERE path[1] IN (SELECT post_id FROM posts WHERE thread_id = $1 AND parent = 0 AND path[1] >
					(SELECT path[1] FROM posts WHERE post_id = $2) ORDER BY post_id ASC LIMIT $3) 
					ORDER BY path ASC, post_id ASC;`
		}

		values = []interface{}{thread.ID, params.Since, params.Limit}
	}

	res := make([]models.Post, 0)

	rows, err = t.db.QueryContext(ctx, query, values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrSuchPostNotFound
		}

		return nil, pkg.ErrWorkDatabase
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		timeTmp := time.Time{}

		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author.Nickname,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&timeTmp)
		if err != nil {
			return nil, err
		}

		post.Thread = thread.ID

		post.Created = timeTmp.Format(time.RFC3339)

		res = append(res, post)
	}

	return res, nil
}
