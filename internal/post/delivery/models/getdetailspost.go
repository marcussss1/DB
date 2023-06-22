package models

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"project/internal/models"
	"project/internal/pkg"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getdetailspost.go

type PostGetDetailsRequest struct {
	ID      int64
	Related []string
}

func NewPostGetDetailsRequest() *PostGetDetailsRequest {
	return &PostGetDetailsRequest{}
}

func (req *PostGetDetailsRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	param := vars["id"]

	value, _ := strconv.Atoi(param)

	req.ID = int64(value)

	param = r.URL.Query().Get("related")

	req.Related = strings.Split(param, ",")

	return nil
}

func (req *PostGetDetailsRequest) GetPost() *models.Post {
	return &models.Post{
		ID: req.ID,
	}
}

func (req *PostGetDetailsRequest) GetParams() *pkg.PostDetailsParams {
	return &pkg.PostDetailsParams{
		Related: req.Related,
	}
}

//easyjson:json
type PostGetDetailsAuthorResponse struct {
	Nickname string `json:"nickname,omitempty"`
	FullName string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
}

//easyjson:json
type PostGetDetailsPostResponse struct {
	ID       int64  `json:"id,omitempty"`
	Parent   int64  `json:"parent,omitempty"`
	Author   string `json:"author,omitempty"`
	Message  string `json:"message,omitempty"`
	IsEdited bool   `json:"isEdited,omitempty"`
	Forum    string `json:"forum,omitempty"`
	Thread   int64  `json:"thread,omitempty"`
	Created  string `json:"created,omitempty"`
}

//easyjson:json
type PostGetDetailsThreadResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Slug    string `json:"slug"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int64  `json:"votes"`
}

//easyjson:json
type PostGetDetailsForumResponse struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts"`
	Threads int64  `json:"threads"`
}

//easyjson:json
type PostGetDetailsResponse struct {
	Post   *PostGetDetailsPostResponse   `json:"post"`
	Thread *PostGetDetailsThreadResponse `json:"thread"`
	Author *PostGetDetailsAuthorResponse `json:"author"`
	Forum  *PostGetDetailsForumResponse  `json:"forum"`
}

func NewPostDetailsResponse(postDetails *models.PostDetails) *PostGetDetailsResponse {
	res := &PostGetDetailsResponse{}

	if postDetails.Post.ID != 0 {
		post := PostGetDetailsPostResponse{
			ID:       postDetails.Post.ID,
			Parent:   postDetails.Post.Parent,
			Author:   postDetails.Post.Author.Nickname,
			Forum:    postDetails.Post.Forum,
			Thread:   postDetails.Post.Thread,
			Message:  postDetails.Post.Message,
			Created:  postDetails.Post.Created,
			IsEdited: postDetails.Post.IsEdited,
		}

		res.Post = &post
	}

	if postDetails.Author.Nickname != "" {
		author := PostGetDetailsAuthorResponse{
			Nickname: postDetails.Author.Nickname,
			FullName: postDetails.Author.FullName,
			About:    postDetails.Author.About,
			Email:    postDetails.Author.Email,
		}

		res.Author = &author
	}

	if postDetails.Thread.ID != 0 {
		thread := PostGetDetailsThreadResponse{
			ID:      postDetails.Thread.ID,
			Title:   postDetails.Thread.Title,
			Author:  postDetails.Thread.Author,
			Forum:   postDetails.Thread.Forum,
			Slug:    postDetails.Thread.Slug,
			Message: postDetails.Thread.Message,
			Created: postDetails.Thread.Created,
			Votes:   postDetails.Thread.Votes,
		}

		res.Thread = &thread
	}

	if postDetails.Forum.User != "" {
		forum := PostGetDetailsForumResponse{
			Title:   postDetails.Forum.Title,
			User:    postDetails.Forum.User,
			Slug:    postDetails.Forum.Slug,
			Posts:   postDetails.Forum.Posts,
			Threads: postDetails.Forum.Threads,
		}

		res.Forum = &forum
	}

	return res
}
