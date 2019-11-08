package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PhilLar/go-chi_example/handlers"
	"github.com/PhilLar/go-chi_example/models"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"

	customMiddleware "github.com/PhilLar/go-chi_example/middleware"
	"github.com/PhilLar/go-chi_example/newsfeed"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var port int
var db string

func init() {
	defPort := 3333
	var defDB string
	if portVar, ok := os.LookupEnv("PORT"); ok {
		if portValue, err := strconv.Atoi(portVar); err == nil {
			defPort = portValue
		}
	}
	if dbVar, ok := os.LookupEnv("DATABASE_URL"); ok {
		defDB = dbVar
	}
	flag.IntVar(&port, "port", defPort, "port to listen on")
	flag.StringVar(&db, "db", defDB, "database to connect to")
}

func main() {

	dbPsql, err := models.NewDB(db, "")
	if err != nil {
		log.Panic(err)
	}
	defer dbPsql.Close()

	env := &handlers.Env{Store: &models.Store{DB: dbPsql}}

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	// You can use [middleware] along the whole app work
	//r.Use(middleware.AllowContentType("application/json"))

	feed := newsfeed.New()
	feed.Add(newsfeed.Item{
		Title: "Hello",
		Post:  "World",
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})

	r.Get("/dbTest", func(w http.ResponseWriter, r *http.Request) {
		err = dbPsql.Ping()
		if err != nil {
			panic(err)
		}
		log.Print("DB OK!")
		w.Write([]byte("DB_TABLE pets CREATED!"))
	})

	r.Get("/newsfeed", handlers.NewsfeedGet(feed))

	// Or you can use [middleware] [With] some special routes
	r.With(middleware.AllowContentType("application/sql")).Post("/newsfeed", handlers.NewsfeedPost(feed))
	r.With(middleware.AllowContentType("application/json")).Post("/newsfeedRIGHT", handlers.NewsfeedPost(feed))

	r.Route("/say", func(r chi.Router) {
		r.Get("/{name}", handlers.RequestSay)
		r.Get("/", handlers.RequestSay)
	})

	r.Route("/pets", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello! glad to see you in pets root!"))
		})
		r.Put("/generate/{amount}", env.PutGeneratePetsHandler())
		r.With(customMiddleware.PetMiddleware).Route("/{kind}/{name}", func(r chi.Router) {
			r.Get("/", handlers.GetPetHandler)
			r.Put("/", env.PutPetHandler())
		})
		r.Get("/{name}", handlers.RequestSay)
		r.Route("/get", func(r chi.Router) {
			r.Get("/all", env.ListPetsHandler())
			r.Get("/filter/{underage}/{overage}/{kind}/{first_letter}", env.FilterPetsHandler())
		})
		r.Route("/remove", func(r chi.Router) {
			r.Put("/all", env.RemoveAllPetsHandler())
		})
	})

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}
