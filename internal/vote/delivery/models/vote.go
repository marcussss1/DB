package models

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"project/internal/models"
	"project/internal/pkg"
)

//go:generate easyjson -all -disallow_unknown_fields -omit_empty vote.go

type VoteRequest struct {
	SlugOrID string
	Nickname string `json:"nickname"`
	Voice    int64  `json:"voice"`
}

func NewVoteRequest() *VoteRequest {
	return &VoteRequest{}
}

func (req *VoteRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	req.SlugOrID = vars["slug_or_id"]

	body, _ := io.ReadAll(r.Body)

	easyjson.Unmarshal(body, req)

	return nil
}

func (req *VoteRequest) GetThread() *models.Thread {
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

func (req *VoteRequest) GetParams() *pkg.VoteParams {
	return &pkg.VoteParams{
		Nickname: req.Nickname,
		Voice:    req.Voice,
	}
}

type VoteResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Slug    string `json:"slug"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int64  `json:"votes"`
}

func NewVoteResponse(thread *models.Thread) *VoteResponse {
	return &VoteResponse{
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
