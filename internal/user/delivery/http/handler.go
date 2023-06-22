package http

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"project/internal/pkg"
	"project/internal/user/delivery/models"
	"project/internal/user/usecase"
)

type UserHandler struct {
	userUsecase usecase.UserService
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewUserCreateRequest()

	err := request.Bind(r)
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	users, err := h.userUsecase.CreateUser(r.Context(), request.GetUser())
	if err != nil {
		if errors.Is(errors.Cause(err), pkg.ErrSuchUserExist) {
			response := models.NewUsersCreateResponse(users)

			pkg.Response(r.Context(), w, http.StatusConflict, response)

			return
		}

		pkg.DefaultHandlerHTTPError(r.Context(), w, err)

		return
	}

	response := models.NewUserCreateResponse(&users[0])

	pkg.Response(r.Context(), w, http.StatusCreated, response)
}

func (h *UserHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewProfileGetRequest()

	request.Bind(r)

	user, err := h.userUsecase.GetProfile(r.Context(), request.GetUser())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewProfileGetResponse(&user)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func (h *UserHandler) UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewProfileUpdateRequest()

	request.Bind(r)

	user, err := h.userUsecase.UpdateProfile(r.Context(), request.GetUser())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewProfileUpdateResponse(&user)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func NewUserHandler(userUsecase usecase.UserService, r *mux.Router) *UserHandler {
	h := &UserHandler{userUsecase: userUsecase}
	return h
}
