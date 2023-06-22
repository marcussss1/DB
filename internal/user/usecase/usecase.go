package usecase

import (
	"context"
	"project/internal/models"
	"project/internal/pkg"
	"project/internal/user/repository"

	"github.com/pkg/errors"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) ([]models.User, error)
	GetProfile(ctx context.Context, user *models.User) (models.User, error)
	UpdateProfile(ctx context.Context, user *models.User) (models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		userRepo: r,
	}
}

func (u userService) CreateUser(ctx context.Context, user *models.User) ([]models.User, error) {
	res, err := u.userRepo.GetUserByEmailOrNickname(ctx, user)
	if err == nil {
		return res, errors.Wrap(pkg.ErrSuchUserExist, "CreateUser")
	}

	userNew, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "CreateUser")
	}

	resOne := []models.User{userNew}

	return resOne, nil
}

func (u userService) GetProfile(ctx context.Context, user *models.User) (models.User, error) {
	res, err := u.userRepo.GetUserByNickname(ctx, user)
	if err != nil {
		return models.User{}, errors.Wrap(err, "GetProfile")
	}

	return res, nil
}

func (u userService) UpdateProfile(ctx context.Context, user *models.User) (models.User, error) {
	_, err := u.userRepo.GetUserByNickname(ctx, user)
	if err != nil {
		return models.User{}, errors.Wrap(err, "UpdateUser")
	}

	exist, _ := u.userRepo.CheckFreeEmail(ctx, user)
	if exist {
		return models.User{}, errors.Wrap(pkg.ErrUpdateUserDataConflict, "UpdateUser")
	}

	resUpdate, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return models.User{}, errors.Wrap(err, "UpdateUser")
	}

	return resUpdate, nil
}
