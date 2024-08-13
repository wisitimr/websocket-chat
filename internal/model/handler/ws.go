package handler

import (
	"net/http"
)

type WsHandler interface {
	WebSocketConnect(w http.ResponseWriter, r *http.Request)
}
