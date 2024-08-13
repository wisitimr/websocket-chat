package service

import (
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"

	"github.com/sirupsen/logrus"
)

type verifyContactService struct {
	redisRepo mRepo.RedisRepository
	logger    *logrus.Logger
}

func InitVerifyContactService(repo mRepo.Repository, logger *logrus.Logger) mService.VerifyContactService {
	return &verifyContactService{
		logger:    logger,
		redisRepo: repo.Redis,
	}
}

func (s verifyContactService) VerifyContact(username string) error {
	if !s.redisRepo.IsUserExist(username) {
		return nil
	}

	return nil
}
