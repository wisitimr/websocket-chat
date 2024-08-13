package handler

import "github.com/go-chi/jwtauth/v5"

type Handler struct {
	AuthToken     *jwtauth.JWTAuth
	WebSocket     WsHandler
	Register      RegisterHandler
	Login         LoginHandler
	VerifyContact VerifyContactHandler
	ChatHistory   ChatHistoryHandler
	ContactList   ContactListHandler
}
