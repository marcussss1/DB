package models

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"project/internal/models"
)

//go:generate easyjson -all -disallow_unknown_fields -omit_empty updatepost.go

type PostUpdateRequest struct {
	ID      int64
	Message string `json:"message"`
}

func NewPostUpdateRequest() *PostUpdateRequest {
	return &PostUpdateRequest{}
}

func (req *PostUpdateRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	param := vars["id"]

	value, _ := strconv.Atoi(param)

	req.ID = int64(value)

	body, _ := io.ReadAll(r.Body)

	easyjson.Unmarshal(body, req)

	return nil
}

func (req *PostUpdateRequest) GetPost() *models.Post {
	return &models.Post{
		ID:      req.ID,
		Message: req.Message,
	}
}

type PostUpdateResponse struct {
	ID       int64  `json:"id"`
	Parent   int64  `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int64  `json:"thread"`
	Created  string `json:"created"`
}

func NewPostUpdateResponse(post *models.Post) *PostUpdateResponse {
	return &PostUpdateResponse{
		ID:       post.ID,
		Parent:   post.Parent,
		Author:   post.Author.Nickname,
		Forum:    post.Forum,
		Thread:   post.Thread,
		Message:  post.Message,
		Created:  post.Created,
		IsEdited: post.IsEdited,
	}
}
