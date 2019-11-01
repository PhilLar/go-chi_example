package main

import (
	"net/http"

	"github.com/PhilLar/go_chi_sample/newsfeed"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	_ = newsfeed.New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	http.ListenAndServe(":3333", r)
}
