package models

type Post struct {
	ID       int64
	Parent   int64
	Author   User
	Message  string
	IsEdited bool
	Forum    string
	Thread   int64
	Created  string
}
