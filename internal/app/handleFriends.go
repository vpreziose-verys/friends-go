package app

import (
	"net/http"
)

func (f *Friends) GetFriends(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetFriends"))
}
