package handler

import (
	"net/http"
)

func (h *Handler) GetFriends(w http.ResponseWriter, r *http.Request) {
	h.base.LogTest()
	w.Write([]byte("GetFriends!"))
}
