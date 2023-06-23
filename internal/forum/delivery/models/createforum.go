package models

import (
	"io"
	"net/http"

	"github.com/mailru/easyjson"

	"project/internal/models"
)

//go:generate easyjson -all -disallow_unknown_fields -omit_empty createforum.go

type ForumCreateRequest struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

func NewForumCreateRequest() *ForumCreateRequest {
	return &ForumCreateRequest{}
}

func (req *ForumCreateRequest) Bind(r *http.Request) error {
	// if r.Header.Get("Content-Type") == "" {
	//	return pkg.ErrContentTypeUndefined
	// }
	//
	// if r.Header.Get("Content-Type") != pkg.ContentTypeJSON {
	//	return pkg.ErrUnsupportedMediaType
	// }

	body, _ := io.ReadAll(r.Body)
	// if err != nil {
	//	return pkg.ErrBadBodyRequest
	// }
	// defer func() {
	//	err = r.Body.Close()
	//	if err != nil {
	//		logrus.Error(err)
	//	}
	// }()

	// if len(body) == 0 {
	//	return pkg.ErrEmptyBody
	// }

	easyjson.Unmarshal(body, req)
	// err = easyjson.Unmarshal(body, req)
	// if err != nil {
	//	return pkg.ErrJSONUnexpectedEnd
	// }

	return nil
}

func (req *ForumCreateRequest) GetForum() *models.Forum {
	return &models.Forum{
		Title: req.Title,
		User:  req.User,
		Slug:  req.Slug,
	}
}

type ForumCreateResponse struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int64  `json:"threads,omitempty"`
}

func NewForumCreateResponse(forum *models.Forum) *ForumCreateResponse {
	return &ForumCreateResponse{
		Title:   forum.Title,
		User:    forum.User,
		Slug:    forum.Slug,
		Posts:   forum.Posts,
		Threads: forum.Threads,
	}
}
