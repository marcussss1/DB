package models

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"project/internal/models"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getdetailsthread.go

type ThreadGetDetailsRequest struct {
	SlugOrID string
}

func NewThreadGetDetailsRequest() *ThreadGetDetailsRequest {
	return &ThreadGetDetailsRequest{}
}

func (req *ThreadGetDetailsRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	req.SlugOrID = vars["slug_or_id"]

	return nil
}

func (req *ThreadGetDetailsRequest) GetThread() *models.Thread {
	id, err := strconv.Atoi(req.SlugOrID)
	if err == nil {
		return &models.Thread{
			ID: int64(id),
		}
	}

	return &models.Thread{
		Slug: req.SlugOrID,
	}
}

//easyjson:json
type ThreadGetDetailsResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Slug    string `json:"slug"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int64  `json:"votes"`
}

func NewThreadGetDetailsResponse(thread *models.Thread) *ThreadGetDetailsResponse {
	return &ThreadGetDetailsResponse{
		ID:      thread.ID,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Slug:    thread.Slug,
		Message: thread.Message,
		Created: thread.Created,
		Votes:   thread.Votes,
	}
}
