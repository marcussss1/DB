package pkg

type GetThreadsParams struct {
	Limit int64
	Since string
	Desc  bool
}

type GetUsersParams struct {
	Limit int64
	Since string
	Desc  bool
}

type GetPostsParams struct {
	Limit int64
	Since int64
	Desc  bool
	Sort  string
}

type VoteParams struct {
	Nickname string
	Voice    int64
}

type PostDetailsParams struct {
	Related []string
}
