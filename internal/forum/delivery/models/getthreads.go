package models

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"project/internal/models"
	"project/internal/pkg"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getthreads.go

type ForumGetThreadsRequest struct {
	Slug  string
	Limit int64
	Since string
	Desc  bool
}

func NewForumGetThreadsRequest() *ForumGetThreadsRequest {
	return &ForumGetThreadsRequest{}
}

func (req *ForumGetThreadsRequest) Bind(r *http.Request) error {
	// if r.Header.Get("Content-Type") != "" {
	//	return pkg.ErrUnsupportedMediaType
	// }

	vars := mux.Vars(r)

	req.Slug = vars["slug"]

	param := ""

	param = r.URL.Query().Get("limit")

	if param != "" {
		value, _ := strconv.Atoi(param)
		// if err != nil {
		//	return pkg.ErrConvertQueryType
		// }

		req.Limit = int64(value)
	} else {
		req.Limit = 100
	}

	// if err != nil {
	//	return pkg.ErrConvertQueryType
	// }

	req.Since = r.FormValue("since")
	// if req.Since == "" {
	//	return pkg.ErrBadRequestParamsEmptyRequiredFields
	// }

	param = r.FormValue("desc")
	// if param == "" {
	//	return pkg.ErrBadRequestParamsEmptyRequiredFields
	// } else if param == "true" {
	//	req.Desc = true
	// } else if param == "false" {
	//	req.Desc = false
	// } else {
	//	return pkg.ErrBadRequestParams
	// }

	if param == "true" {
		req.Desc = true
	}

	return nil
}

func (req *ForumGetThreadsRequest) GetForum() *models.Forum {
	return &models.Forum{
		Slug: req.Slug,
	}
}

func (req *ForumGetThreadsRequest) GetParams() *pkg.GetThreadsParams {
	return &pkg.GetThreadsParams{
		Limit: req.Limit,
		Since: req.Since,
		Desc:  req.Desc,
	}
}

//easyjson:json
type ForumGetThreadsResponse struct {
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
type ThreadsList []ForumGetThreadsResponse

func NewForumGetThreadsResponse(threads []*models.Thread) ThreadsList {
	res := make([]ForumGetThreadsResponse, len(threads))

	for idx, value := range threads {
		res[idx] = ForumGetThreadsResponse{
			ID:      value.ID,
			Title:   value.Title,
			Author:  value.Author,
			Forum:   value.Forum,
			Slug:    value.Slug,
			Message: value.Message,
			Created: value.Created,
			Votes:   value.Votes,
		}
	}

	return res
}
