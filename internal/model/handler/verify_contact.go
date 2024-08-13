package handler

import (
	"net/http"
)

type VerifyContactHandler interface {
	VerifyContact(w http.ResponseWriter, r *http.Request)
}
