package service

import (
	"websocket-chat/internal/model/user"
)

type RegisterService interface {
	NewUser(user user.Request) error
}
