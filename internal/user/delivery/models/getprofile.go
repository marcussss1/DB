package models

import (
	"net/http"

	"github.com/gorilla/mux"

	"project/internal/models"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty getprofile.go

type ProfileGetRequest struct {
	Nickname string
}

func NewProfileGetRequest() *ProfileGetRequest {
	return &ProfileGetRequest{}
}

func (req *ProfileGetRequest) Bind(r *http.Request) error {
	vars := mux.Vars(r)

	req.Nickname = vars["nickname"]

	return nil
}

func (req *ProfileGetRequest) GetUser() *models.User {
	return &models.User{
		Nickname: req.Nickname,
	}
}

//easyjson:json
type ProfileGetResponse struct {
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

func NewProfileGetResponse(user *models.User) *ProfileGetResponse {
	return &ProfileGetResponse{
		Nickname: user.Nickname,
		FullName: user.FullName,
		About:    user.About,
		Email:    user.Email,
	}
}
