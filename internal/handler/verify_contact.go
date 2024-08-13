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

type verifyContactHandler struct {
	verifyContactService mService.VerifyContactService
	logger               *logrus.Logger
	mRes.ResponseDto
}

func InitVerifyContactHandler(service mService.Service, logger *logrus.Logger) mHandler.VerifyContactHandler {
	return verifyContactHandler{
		verifyContactService: service.VerifyContact,
		logger:               logger,
	}
}

func (h verifyContactHandler) VerifyContact(w http.ResponseWriter, r *http.Request) {
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

	err = h.verifyContactService.VerifyContact(u.Username)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}

	h.Respond(w, r, nil, http.StatusOK)
}
