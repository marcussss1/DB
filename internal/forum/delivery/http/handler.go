package http

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"project/internal/forum/delivery/models"
	"project/internal/forum/usecase"
	"project/internal/pkg"
)

type ForumHandler struct {
	forumUsecase usecase.ForumService
}

func (h *ForumHandler) CreateForumHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewForumCreateRequest()

	request.Bind(r)

	forum, err := h.forumUsecase.CreateForum(r.Context(), request.GetForum())
	if err != nil {
		if errors.Is(errors.Cause(err), pkg.ErrSuchForumExist) {
			response := models.NewForumCreateResponse(forum)

			pkg.Response(r.Context(), w, http.StatusConflict, response)

			return
		}

		pkg.DefaultHandlerHTTPError(r.Context(), w, err)

		return
	}

	response := models.NewForumCreateResponse(forum)

	pkg.Response(r.Context(), w, http.StatusCreated, response)
}

func (h *ForumHandler) GetForumHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewForumGetDetailsRequest()

	request.Bind(r)

	forum, err := h.forumUsecase.GetDetailsForum(r.Context(), request.GetForum())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewForumGetDetailsResponse(forum)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func (h *ForumHandler) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	request := models.NewForumGetThreadsRequest()

	request.Bind(r)

	threads, err := h.forumUsecase.GetThreads(r.Context(), request.GetForum(), request.GetParams())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewForumGetThreadsResponse(threads)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func (h *ForumHandler) GetForumUsersHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewForumGetUsersRequest()

	request.Bind(r)

	users, err := h.forumUsecase.GetUsers(r.Context(), request.GetForum(), request.GetParams())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewForumGetUsersResponse(users)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func NewForumHandler(forumUsecase usecase.ForumService, r *mux.Router) *ForumHandler {
	h := &ForumHandler{forumUsecase: forumUsecase}
	return h
}
