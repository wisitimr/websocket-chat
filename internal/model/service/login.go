package service

import (
	"websocket-chat/internal/model/user"
)

type LoginService interface {
	Login(user user.Request) error
}
