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
		kind := chi.URLParam(r, "kind")


		if _, ok := allowedPets[kind]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("forbidden pet of kind: %s!\n", kind)))
			return
		}

		next.ServeHTTP(w, r)
	})
}