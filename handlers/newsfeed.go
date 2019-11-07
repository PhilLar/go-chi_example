package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/PhilLar/go-chi_example/newsfeed"
	"github.com/go-chi/chi"
	"net/http"
)

func NewsfeedGet(feed newsfeed.Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := feed.GetAll()
		json.NewEncoder(w).Encode(items)
	}
}

func RequestSay(w http.ResponseWriter, r *http.Request) {
	val := chi.URLParam(r, "name")
	if val != "" {
		fmt.Fprintf(w, "Hello %s!", val)
	} else {
		fmt.Fprintf(w, "Hello ... you.")
	}
}

func NewsfeedPost(feed newsfeed.Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := map[string]string{}
		json.NewDecoder(r.Body).Decode(&request)

		feed.Add(newsfeed.Item{
			Title:	request["title"],
			Post:	request["post"],
		})
		w.Write([]byte("alright!\n"))
	}
}