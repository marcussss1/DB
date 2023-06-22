package usecase

import (
	"context"

	"github.com/pkg/errors"

	repoForum "project/internal/forum/repository"
	"project/internal/models"
	"project/internal/pkg"
	repoUser "project/internal/user/repository"
)

type ForumService interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	GetDetailsForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	GetThreads(ctx context.Context, forum *models.Forum, params *pkg.GetThreadsParams) ([]*models.Thread, error)
	GetUsers(ctx context.Context, forum *models.Forum, params *pkg.GetUsersParams) ([]*models.User, error)
}

type forumService struct {
	forumRepo repoForum.ForumRepository
	userRepo  repoUser.UserRepository
}

func NewForumService(rf repoForum.ForumRepository, ru repoUser.UserRepository) ForumService {
	return &forumService{
		forumRepo: rf,
		userRepo:  ru,
	}
}

func (f forumService) CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	res, err := f.forumRepo.GetDetailsForumBySlug(ctx, forum)
	if err == nil {
		return res, errors.Wrap(pkg.ErrSuchForumExist, "CreateForum")
	}

	user, err := f.userRepo.GetUserByNickname(ctx, &models.User{Nickname: forum.User})
	if err != nil {
		return res, errors.Wrap(err, "CreateForum")
	}

	forum.User = user.Nickname

	res, err = f.forumRepo.CreateForum(ctx, forum)
	if err != nil {
		_, err = f.userRepo.GetUserByNickname(ctx, &models.User{Nickname: forum.User})
		if err != nil {
			return nil, errors.Wrap(err, "CreateForum")
		}
	}

	return res, nil
}

func (f forumService) GetDetailsForum(ctx context.Context, forum *models.Forum) (*models.Forum, error) {
	res, err := f.forumRepo.GetDetailsForumBySlug(ctx, forum)
	if err != nil {
		return nil, errors.Wrap(err, "GetDetailsForumBySlug")
	}

	return res, nil
}

func (f forumService) GetThreads(ctx context.Context, forum *models.Forum, params *pkg.GetThreadsParams) ([]*models.Thread, error) {
	exist, _ := f.forumRepo.CheckExistForum(ctx, forum)
	if !exist {
		return nil, errors.Wrap(pkg.ErrSuchForumNotFound, "GetThreads")
	}

	res, err := f.forumRepo.GetThreads(ctx, forum, params)
	if err != nil {
		return nil, errors.Wrap(err, "GetThreads")
	}

	return res, nil
}

func (f forumService) GetUsers(ctx context.Context, forum *models.Forum, params *pkg.GetUsersParams) ([]*models.User, error) {
	exist, _ := f.forumRepo.CheckExistForum(ctx, forum)
	if !exist {
		return nil, errors.Wrap(pkg.ErrSuchForumNotFound, "GetThreads")
	}

	res, err := f.forumRepo.GetUsers(ctx, forum, params)
	if err != nil {
		return nil, errors.Wrap(err, "GetUsers")
	}

	return res, nil
}
