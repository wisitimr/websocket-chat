package handler

import (
	"net/http"
)

type RegisterHandler interface {
	NewUser(w http.ResponseWriter, r *http.Request)
}
