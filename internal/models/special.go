package models

type StatusService struct {
	User   int64
	Forum  int64
	Thread int64
	Post   int64
}

type PostDetails struct {
	Post   Post
	Author User
	Thread Thread
	Forum  Forum
}
