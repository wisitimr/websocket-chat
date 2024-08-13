package service

import "websocket-chat/internal/model"

type ContactListService interface {
	ContactList(username1, username2, fromTS, toTS string) ([]model.Chat, error)
}
