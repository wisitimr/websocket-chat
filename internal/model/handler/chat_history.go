package handler

import (
	"net/http"
)

type ChatHistoryHandler interface {
	ChatHistory(w http.ResponseWriter, r *http.Request)
}
