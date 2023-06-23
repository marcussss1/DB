package models

import (
	"net/http"

	"github.com/gorilla/mux"

	"project/internal/models"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getdetailsforum.go

type ForumGetDetailsRequest struct {
	Slug string
}

func NewForumGetDetailsRequest() *ForumGetDetailsRequest {
	return &ForumGetDetailsRequest{}
}

func (req *ForumGetDetailsRequest) Bind(r *http.Request) error {
	// if r.Header.Get("Content-Type") != "" {
	//	return pkg.ErrUnsupportedMediaType
	// }

	vars := mux.Vars(r)

	req.Slug = vars["slug"]

	return nil
}

func (req *ForumGetDetailsRequest) GetForum() *models.Forum {
	return &models.Forum{
		Slug: req.Slug,
	}
}

//easyjson:json
type ForumGetDetailsResponse struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int64  `json:"threads,omitempty"`
}

func NewForumGetDetailsResponse(forum *models.Forum) *ForumGetDetailsResponse {
	return &ForumGetDetailsResponse{
		Title:   forum.Title,
		User:    forum.User,
		Slug:    forum.Slug,
		Posts:   forum.Posts,
		Threads: forum.Threads,
	}
}
