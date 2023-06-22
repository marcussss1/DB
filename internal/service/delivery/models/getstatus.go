package models

import "project/internal/models"

//go:generate easyjson -all -disallow_unknown_fields -omit_empty getstatus.go

type ServiceGetStatusResponse struct {
	User   int64 `json:"user"`
	Forum  int64 `json:"forum"`
	Thread int64 `json:"thread"`
	Post   int64 `json:"post"`
}

func NewServiceGetStatusResponse(service *models.StatusService) *ServiceGetStatusResponse {
	return &ServiceGetStatusResponse{
		User:   service.User,
		Forum:  service.Forum,
		Thread: service.Thread,
		Post:   service.Post,
	}
}
