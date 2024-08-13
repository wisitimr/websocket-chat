package service

import (
	"websocket-chat/internal/model"
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"

	"github.com/sirupsen/logrus"
)

type contactListService struct {
	redisRepo mRepo.RedisRepository
	logger    *logrus.Logger
}

func InitContactListService(repo mRepo.Repository, logger *logrus.Logger) mService.ContactListService {
	return &contactListService{
		logger:    logger,
		redisRepo: repo.Redis,
	}
}

func (s contactListService) ContactList(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	if !s.redisRepo.IsUserExist(username1) || !s.redisRepo.IsUserExist(username2) {
		return nil, nil
	}

	chats, err := s.redisRepo.FetchChatBetween(username1, username2, fromTS, toTS)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
