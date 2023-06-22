package usecase

import (
	"context"

	"github.com/pkg/errors"

	"project/internal/models"
	"project/internal/pkg"
	threadRepo "project/internal/thread/repository"
	userRepo "project/internal/user/repository"
	voteRepo "project/internal/vote/repository"
)

type VoteService interface {
	Vote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) (models.Thread, error)
}

type voteService struct {
	voteRepo   voteRepo.VoteRepository
	threadRepo threadRepo.ThreadRepository
	userRepo   userRepo.UserRepository
}

func NewVoteService(vr voteRepo.VoteRepository, tr threadRepo.ThreadRepository, ur userRepo.UserRepository) VoteService {
	return &voteService{
		voteRepo:   vr,
		threadRepo: tr,
		userRepo:   ur,
	}
}

func (v voteService) Vote(ctx context.Context, thread *models.Thread, params *pkg.VoteParams) (models.Thread, error) {
	var err error

	resThread := models.Thread{}

	// CheckAndGetThread
	if thread.Slug != "" {
		resThread, err = v.threadRepo.GetDetailsThreadBySlug(ctx, thread)
	} else {
		resThread, err = v.threadRepo.GetDetailsThreadByID(ctx, thread)
	}
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "Vote")
	}

	// CheckUser
	resUser, err := v.userRepo.GetUserByNickname(ctx, &models.User{Nickname: params.Nickname})
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "Vote")
	}
	params.Nickname = resUser.Nickname

	// CheckVote
	exist, err := v.voteRepo.CheckExistVote(ctx, &resThread, params)
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "Vote CheckExistVote")
	}

	if exist {
		err = v.voteRepo.UpdateVote(ctx, &resThread, params)
	} else {
		err = v.voteRepo.CreateVote(ctx, &resThread, params)
	}
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "Vote exist")
	}

	threadUPD, err := v.threadRepo.GetDetailsThreadByID(ctx, &resThread)
	if err != nil {
		return models.Thread{}, errors.Wrap(err, "Vote GetDetailsThreadByID")
	}

	return threadUPD, nil
}
