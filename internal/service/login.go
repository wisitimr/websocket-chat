package service

import (
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"
	"websocket-chat/internal/model/user"

	"github.com/sirupsen/logrus"
)

type loginService struct {
	redisRepo mRepo.RedisRepository
	logger    *logrus.Logger
}

func InitLoginService(repo mRepo.Repository, logger *logrus.Logger) mService.LoginService {
	return &loginService{
		logger:    logger,
		redisRepo: repo.Redis,
	}
}

func (s loginService) Login(u user.Request) error {
	err := s.redisRepo.IsUserAuthentic(u.Username, u.Password)
	if err != nil {
		return err
	}

	return nil
}
