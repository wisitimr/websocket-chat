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

type loginHandler struct {
	loginService mService.LoginService
	logger       *logrus.Logger
	mRes.ResponseDto
}

func InitLoginHandler(service mService.Service, logger *logrus.Logger) mHandler.LoginHandler {
	return loginHandler{
		loginService: service.Login,
		logger:       logger,
	}
}

func (h loginHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	err = h.loginService.Login(u)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}

	h.Respond(w, r, nil, http.StatusOK)
}
