package middleware

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func PetMiddleware(next http.Handler) http.Handler {
	allowedPets := map[string]struct{}{
		"cat": struct{}{},
		"dog": struct{}{},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pet := chi.URLParam(r, "pet")

		if _, ok := allowedPets[pet]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("forbidden pet: %s!\n", pet)))
			return
		}

		next.ServeHTTP(w, r)
	})
}