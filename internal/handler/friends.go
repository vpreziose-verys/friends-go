package handler

import (
	"net/http"
)

func (h *Handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetFriends Yo!"))
}
