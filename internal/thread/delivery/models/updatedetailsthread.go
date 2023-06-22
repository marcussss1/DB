package models

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"project/internal/models"
)

//go:generate easyjson -all -disallow_unknown_fields -omit_empty updatedetailsthread.go

type ThreadUpdateDetailsRequest struct {
	SlugOrID string
	Title    string `json:"title"`
	Message  string `json:"message"`
}

func NewThreadUpdateDetailsRequest() *ThreadUpdateDetailsRequest {
	return &ThreadUpdateDetailsRequest{}
}

func (req *ThreadUpdateDetailsRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	req.SlugOrID = vars["slug_or_id"]

	body, _ := io.ReadAll(r.Body)

	easyjson.Unmarshal(body, req)

	return nil
}

func (req *ThreadUpdateDetailsRequest) GetThread() *models.Thread {
	id, err := strconv.Atoi(req.SlugOrID)
	if err == nil {
		return &models.Thread{
			ID:      int64(id),
			Message: req.Message,
			Title:   req.Title,
		}
	}

	return &models.Thread{
		Slug:    req.SlugOrID,
		Message: req.Message,
		Title:   req.Title,
	}
}

//easyjson:json
type ThreadUpdateDetailsResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Slug    string `json:"slug"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int64  `json:"votes"`
}

func NewThreadUpdateDetailsResponse(thread *models.Thread) *ThreadUpdateDetailsResponse {
	return &ThreadUpdateDetailsResponse{
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
