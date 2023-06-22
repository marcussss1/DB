package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"project/internal/pkg"
	"project/internal/vote/delivery/models"
	"project/internal/vote/usecase"
)

type VoteHandler struct {
	voteUsecase usecase.VoteService
}

func (h *VoteHandler) VoteHandler(w http.ResponseWriter, r *http.Request) {
	request := models.NewVoteRequest()

	request.Bind(r)

	thread, err := h.voteUsecase.Vote(r.Context(), request.GetThread(), request.GetParams())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewVoteResponse(&thread)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func NewVoteHandler(voteUsecase usecase.VoteService, r *mux.Router) *VoteHandler {
	h := &VoteHandler{voteUsecase: voteUsecase}
	return h
}
