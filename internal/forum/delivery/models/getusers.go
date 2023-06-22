package models

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"project/internal/models"
	"project/internal/pkg"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getusers.go

type ForumGetUsersRequest struct {
	Slug  string
	Limit int64
	Since string
	Desc  bool
}

func NewForumGetUsersRequest() *ForumGetUsersRequest {
	return &ForumGetUsersRequest{}
}

func (req *ForumGetUsersRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	req.Slug = vars["slug"]

	param := ""

	param = r.FormValue("limit")
	if param != "" {
		value, _ := strconv.Atoi(param)
		req.Limit = int64(value)
	} else {
		req.Limit = 100
	}

	req.Since = r.FormValue("since")

	param = r.FormValue("desc")

	if param == "true" {
		req.Desc = true
	}

	return nil
}

func (req *ForumGetUsersRequest) GetForum() *models.Forum {
	return &models.Forum{
		Slug: req.Slug,
	}
}

func (req *ForumGetUsersRequest) GetParams() *pkg.GetUsersParams {
	return &pkg.GetUsersParams{
		Limit: req.Limit,
		Since: req.Since,
		Desc:  req.Desc,
	}
}

//easyjson:json
type ForumGetUsersResponse struct {
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

//easyjson:json
type UsersList []ForumGetUsersResponse

func NewForumGetUsersResponse(users []*models.User) UsersList {
	res := make([]ForumGetUsersResponse, len(users))

	for idx, value := range users {
		res[idx] = ForumGetUsersResponse{
			Nickname: value.Nickname,
			FullName: value.FullName,
			About:    value.About,
			Email:    value.Email,
		}
	}

	return res
}
