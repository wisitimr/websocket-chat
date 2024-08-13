package service

import (
	"websocket-chat/internal/model"
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"

	"github.com/sirupsen/logrus"
)

type chatHistoryService struct {
	redisRepo mRepo.RedisRepository
	logger    *logrus.Logger
}

func InitChatHistoryService(repo mRepo.Repository, logger *logrus.Logger) mService.ChatHistoryService {
	return &chatHistoryService{
		logger:    logger,
		redisRepo: repo.Redis,
	}
}

func (s chatHistoryService) ChatHistory(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	if !s.redisRepo.IsUserExist(username1) || !s.redisRepo.IsUserExist(username2) {
		return nil, nil
	}

	chats, err := s.redisRepo.FetchChatBetween(username1, username2, fromTS, toTS)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
