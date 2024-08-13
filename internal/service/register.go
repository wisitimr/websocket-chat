package service

import (
	"log"
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"
	"websocket-chat/internal/model/user"

	"github.com/sirupsen/logrus"
)

type registerService struct {
	redisRepo mRepo.RedisRepository
	logger    *logrus.Logger
}

func InitRegisterService(repo mRepo.Repository, logger *logrus.Logger) mService.RegisterService {
	return &registerService{
		logger:    logger,
		redisRepo: repo.Redis,
	}
}

func (s registerService) NewUser(u user.Request) error {
	err := s.redisRepo.RegisterNewUser(u.Username, u.Password)
	if err != nil {
		log.Println("error while adding new user", err)
		return err
	}

	return nil
}
