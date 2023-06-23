package models

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"project/internal/models"
	"project/internal/pkg"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getposts.go

type ThreadGetPostsRequest struct {
	SlugOrID string
	Limit    int64
	Since    int64
	Desc     bool
	Sort     string
}

func NewThreadGetPostsRequest() *ThreadGetPostsRequest {
	return &ThreadGetPostsRequest{}
}

func (req *ThreadGetPostsRequest) Bind(r *http.Request) error {
	// if r.Header.Get("Content-Type") != "" {
	//	return pkg.ErrUnsupportedMediaType
	// }

	vars := mux.Vars(r)

	req.SlugOrID = vars["slug_or_id"]

	param := ""

	param = r.FormValue("limit")
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

	param = r.FormValue("since")
	// if req.Since == "" {
	//	return pkg.ErrBadRequestParamsEmptyRequiredFields
	// }
	if param != "" {
		value, _ := strconv.Atoi(param)
		// if err != nil {
		//	return pkg.ErrConvertQueryType
		// }

		req.Since = int64(value)
	} else {
		req.Since = -1
	}

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

	req.Sort = r.FormValue("sort")
	if req.Sort == "" {
		req.Sort = "flat"
	}

	return nil
}

func (req *ThreadGetPostsRequest) GetThread() *models.Thread {
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

func (req *ThreadGetPostsRequest) GetParams() *pkg.GetPostsParams {
	return &pkg.GetPostsParams{
		Limit: req.Limit,
		Since: req.Since,
		Desc:  req.Desc,
		Sort:  req.Sort,
	}
}

//easyjson:json
type ThreadGetPostsResponse struct {
	ID       int64  `json:"id"`
	Parent   int64  `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited"`
	Forum    string `json:"forum"`
	Thread   int64  `json:"thread"`
	Created  string `json:"created"`
}

//easyjson:json
type PostsList []ThreadGetPostsResponse

func NewThreadGetPostsResponse(posts []models.Post) PostsList {
	res := make([]ThreadGetPostsResponse, len(posts))

	for idx, value := range posts {
		res[idx] = ThreadGetPostsResponse{
			ID:       value.ID,
			Parent:   value.Parent,
			Author:   value.Author.Nickname,
			Forum:    value.Forum,
			Thread:   value.Thread,
			Message:  value.Message,
			Created:  value.Created,
			IsEdited: value.IsEdited,
		}
	}

	return res
}
