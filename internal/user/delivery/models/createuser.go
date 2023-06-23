package models

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"project/internal/models"
)

//go:generate easyjson -disallow_unknown_fields -omit_empty createuser.go

//easyjson:json
type UserCreateRequest struct {
	Nickname string
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

func NewUserCreateRequest() *UserCreateRequest {
	return &UserCreateRequest{}
}

func (req *UserCreateRequest) Bind(r *http.Request) error {
	// if r.Header.Get("Content-Type") == "" {
	//	return pkg.ErrContentTypeUndefined
	// }
	//
	// if r.Header.Get("Content-Type") != pkg.ContentTypeJSON {
	//	return pkg.ErrUnsupportedMediaType
	// }

	vars := mux.Vars(r)

	req.Nickname = vars["nickname"]

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

func (req *UserCreateRequest) GetUser() *models.User {
	return &models.User{
		Nickname: req.Nickname,
		FullName: req.FullName,
		About:    req.About,
		Email:    req.Email,
	}
}

//easyjson:json
type UserCreateResponse struct {
	Nickname string `json:"nickname"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

//easyjson:json
type UsersList []UserCreateResponse

func NewUserCreateResponse(user *models.User) *UserCreateResponse {
	return &UserCreateResponse{
		Nickname: user.Nickname,
		FullName: user.FullName,
		About:    user.About,
		Email:    user.Email,
	}
}

func NewUsersCreateResponse(users []models.User) UsersList {
	res := make([]UserCreateResponse, len(users))

	for idx, value := range users {
		res[idx] = UserCreateResponse{
			Nickname: value.Nickname,
			FullName: value.FullName,
			About:    value.About,
			Email:    value.Email,
		}
	}

	return res
}
