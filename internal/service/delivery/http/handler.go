package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"project/internal/pkg"
	"project/internal/service/delivery/models"
	"project/internal/service/usecase"
)

type ServiceHandler struct {
	serviceUsecase usecase.Service
}

func (h *ServiceHandler) ServiceClearHandler(w http.ResponseWriter, r *http.Request) {
	err := h.serviceUsecase.Clear(r.Context())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	pkg.NoBody(w, http.StatusOK)
}

func (h *ServiceHandler) ServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := h.serviceUsecase.GetStatus(r.Context())
	if err != nil {
		pkg.DefaultHandlerHTTPError(r.Context(), w, err)
		return
	}

	response := models.NewServiceGetStatusResponse(status)

	pkg.Response(r.Context(), w, http.StatusOK, response)
}

func NewServiceHandler(serviceUsecase usecase.Service, r *mux.Router) *ServiceHandler {
	h := &ServiceHandler{serviceUsecase: serviceUsecase}
	return h
}
