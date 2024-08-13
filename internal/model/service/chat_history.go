package service

import "websocket-chat/internal/model"

type ChatHistoryService interface {
	ChatHistory(username1, username2, fromTS, toTS string) ([]model.Chat, error)
}
