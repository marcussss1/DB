package pkg

import "database/sql"

const (
	ContentTypeJSON = "application/json"
	BufSizeRequest  = 1024 * 1024 * 1
)

type ContextKeyType string

var SessionKey ContextKeyType = "cookie"

const RequestID = "req-id"

var RequestIDKey ContextKeyType = RequestID

var LoggerKey ContextKeyType = "logger"

var TxInsertOptions = &sql.TxOptions{
	Isolation: sql.LevelDefault,
	ReadOnly:  false,
}

const (
	TypeSortFlat       = "flat"
	TypeSortTree       = "tree"
	TypeSortParentTree = "parent_tree"

	PostDetailForum  = "forum"
	PostDetailThread = "thread"
	PostDetailAuthor = "user"
)
