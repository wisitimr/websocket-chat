package handler

import (
	"encoding/json"
	"net/http"
	mHandler "websocket-chat/internal/model/handler"
	mRes "websocket-chat/internal/model/response"
	mService "websocket-chat/internal/model/service"
	"websocket-chat/internal/model/user"

	"github.com/sirupsen/logrus"
)

type chatHistoryHandler struct {
	chatHistoryService mService.ChatHistoryService
	logger             *logrus.Logger
	mRes.ResponseDto
}

func InitChatHistoryHandler(service mService.Service, logger *logrus.Logger) mHandler.ChatHistoryHandler {
	return chatHistoryHandler{
		chatHistoryService: service.ChatHistory,
		logger:             logger,
	}
}

func (h chatHistoryHandler) ChatHistory(w http.ResponseWriter, r *http.Request) {
	// check if username in userset
	// return error if exist
	// create new user
	// create response for error
	u := user.Request{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}

	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")

	// chat between timerange fromTS toTS
	// where TS is timestamp
	// 0 to positive infinity
	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	chats, err := h.chatHistoryService.ChatHistory(u1, u2, fromTS, toTS)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}

	h.Respond(w, r, chats, http.StatusOK)
}
