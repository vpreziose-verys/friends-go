package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (f *Friends) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(
		//HandleErrors,
		middleware.RequestID,
		middleware.RealIP,
		//Logger,
		//WrapNewRelic(f.nra),
	)

	r.Mount("/", publicRouter(f))
	r.Mount("/public", publicRouter(f))
	//r.Mount("/private", privateRouter(f))

	//r.NotFound(NotFound)
	//r.MethodNotAllowed(MethodNotAllowed)

	return r
}

func publicRouter(f *Friends) http.Handler {
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		r.Route("/friends", func(r chi.Router) {
			r.Get("/", f.GetFriends)
		})
	})

	return r
}
