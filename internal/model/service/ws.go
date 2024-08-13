package service

import "websocket-chat/internal/model"

type WsService interface {
	Receiver(client *model.Client)
}
