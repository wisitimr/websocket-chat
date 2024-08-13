package repository

import "websocket-chat/internal/model"

type RedisRepository interface {
	RegisterNewUser(username, password string) error
	IsUserExist(username string) bool
	IsUserAuthentic(username, password string) error
	CreateChat(c *model.Chat) (string, error)
	CreateFetchChatBetweenIndex()
	FetchChatBetween(username1, username2, fromTS, toTS string) ([]model.Chat, error)
	FetchContactList(username string) ([]model.ContactList, error)
}
