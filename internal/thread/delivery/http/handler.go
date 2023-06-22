package http

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"project/internal/pkg"
	"project/internal/thread/delivery/models"
	"project/internal/thread/usecase"
)

type ThreadHandler struct {
	threadUsecase usecase.ThreadService
}

func (h *ThreadHandler) CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewThreadCreatePostsRequest()

	request.Bind(r)

	posts, err := h.threadUsecase.CreatePosts(r.Context(), request.GetThread(), request.GetPosts())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewThreadCreatePostsResponse(posts)

	pkg.Response(r.Context(), w, http.StatusCreated, response)
}

func (h *ThreadHandler) CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewForumCreateThreadRequest()

	request.Bind(r)

	thread, err := h.threadUsecase.CreateThread(r.Context(), request.GetThread())
	if err != nil {
		if errors.Is(errors.Cause(err), pkg.ErrSuchThreadExist) {
			response := models.NewForumCreateThreadResponse(&thread)

			pkg.Response(r.Context(), w, http.StatusConflict, response)

			return
		}

		pkg.DefaultHandlerHTTPError(r.Context(), w, err)

		return
	}

	response := models.NewForumCreateThreadResponse(&thread)

	pkg.Response(r.Context(), w, http.StatusCreated, response)
}

func (h *ThreadHandler) GetThreadHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewThreadGetDetailsRequest()

	request.Bind(r)

	thread, err := h.threadUsecase.GetDetailsThread(r.Context(), request.GetThread())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewThreadGetDetailsResponse(&thread)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func (h *ThreadHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewThreadGetPostsRequest()

	request.Bind(r)

	posts, err := h.threadUsecase.GetPosts(r.Context(), request.GetThread(), request.GetParams())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewThreadGetPostsResponse(posts)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func (h *ThreadHandler) UpdateThreadHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewThreadUpdateDetailsRequest()

	request.Bind(r)

	thread, err := h.threadUsecase.UpdateThread(r.Context(), request.GetThread())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewThreadUpdateDetailsResponse(&thread)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func NewThreadHandler(threadUsecase usecase.ThreadService, r *mux.Router) *ThreadHandler {
	h := &ThreadHandler{threadUsecase: threadUsecase}
	return h
}
