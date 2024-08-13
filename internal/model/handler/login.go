package handler

import (
	"net/http"
)

type LoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}
