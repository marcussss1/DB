package usecase

import (
	"context"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/service/repository"
)

type Service interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (*models.StatusService, error)
}

type service struct {
	serviceRepo repository.ServiceRepository
}

func NewService(r repository.ServiceRepository) Service {
	return &service{
		serviceRepo: r,
	}
}

func (s service) Clear(ctx context.Context) error {
	err := s.serviceRepo.Clear(ctx)
	if err != nil {
		return errors.Wrap(err, "Clear")
	}

	return err
}

func (s service) GetStatus(ctx context.Context) (*models.StatusService, error) {
	res, err := s.serviceRepo.GetStatus(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "GetStatus")
	}

	return res, nil
}
