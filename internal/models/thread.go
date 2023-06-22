package models

type Thread struct {
	ID      int64
	Title   string
	Author  string
	Forum   string
	Slug    string
	Message string
	Created string
	Votes   int64
}
