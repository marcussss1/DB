package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"project/internal/pkg"
	"project/internal/post/delivery/models"
	"project/internal/post/usecase"
)

type PostHandler struct {
	postUsecase usecase.PostService
}

func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewPostGetDetailsRequest()

	request.Bind(r)

	postDetails, err := h.postUsecase.GetDetailsPost(r.Context(), request.GetPost(), request.GetParams())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewPostDetailsResponse(postDetails)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func (h *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewPostUpdateRequest()

	request.Bind(r)

	post, err := h.postUsecase.UpdatePost(r.Context(), request.GetPost())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewPostUpdateResponse(post)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func NewPostHandler(postUsecase usecase.PostService, r *mux.Router) *PostHandler {
	h := &PostHandler{postUsecase: postUsecase}
	return h
}
