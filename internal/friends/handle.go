package friends

import (
	"net/http"
)

func (f *Friends) GetPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetPing"))
}

func (f *Friends) GetFriends(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetFriends"))
}
