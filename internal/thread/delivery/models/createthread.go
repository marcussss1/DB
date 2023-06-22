package models

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"project/internal/models"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty createthread.go

//easyjson:json
type ForumCreateThreadRequest struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Slug    string `json:"slug"`
}

func NewForumCreateThreadRequest() *ForumCreateThreadRequest {
	return &ForumCreateThreadRequest{}
}

func (req *ForumCreateThreadRequest) Bind(r *http.Request) error {
	body, _ := io.ReadAll(r.Body)

	easyjson.Unmarshal(body, req)

	vars := mux.Vars(r)

	req.Forum = vars["slug"]

	return nil
}

func (req *ForumCreateThreadRequest) GetThread() *models.Thread {
	return &models.Thread{
		Slug:    req.Slug,
		Title:   req.Title,
		Author:  req.Author,
		Message: req.Message,
		Created: req.Created,
		Forum:   req.Forum,
	}
}

//easyjson:json
type ForumCreateThreadResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
	Votes   int64  `json:"votes"`
}

func NewForumCreateThreadResponse(thread *models.Thread) *ForumCreateThreadResponse {
	return &ForumCreateThreadResponse{
		ID:      thread.ID,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Created: thread.Created,
		Votes:   thread.Votes,
		Slug:    thread.Slug,
	}
}
