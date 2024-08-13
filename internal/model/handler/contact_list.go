package handler

import (
	"net/http"
)

type ContactListHandler interface {
	ContactList(w http.ResponseWriter, r *http.Request)
}
