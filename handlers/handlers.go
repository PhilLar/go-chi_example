package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/PhilLar/go-chi_example/models"
	"github.com/PhilLar/go-chi_example/newsfeed"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var randomNameAPI string = "http://names.drycodes.com"

type PetStore interface {
	InsertPet(name, kind string) (int, error)
	ListPets() ([]*models.Pet, error)
}

type Env struct {
	Store       PetStore
}

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

// Pets GET and PUT used with middleware
func GetPetHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	w.Write([]byte(fmt.Sprintf("get pet with name: %s", name)))
}

func (env *Env) PutPetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")
		name := chi.URLParam(r, "name")
		ID, err := env.Store.InsertPet(name, kind)
		if err != nil {
			log.Print(err)
			http.Error(w, "Query to db was not completed", 400)
		}
		err = json.NewEncoder(w).Encode(models.Pet{
			ID:   ID,
			Name: name,
			Kind: kind,
		})
		if err != nil {
			log.Print(err)
			http.Error(w, "Error while json-encoding", 400)
		}
	}
}

func (env *Env) PutGeneratePetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := chi.URLParam(r, "amount")
		amount, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}

		names, err := generateNames(val)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		pets := make([]*models.Pet, 0)
		log.Println("XXXXXXXXXXXXXXXXXXXXX")
		log.Println(names)
		for i:=0; i<amount; i++ {
			name := names[i]
			kind := chooseYourDestiny()
			ID, err := env.Store.InsertPet(name, kind)
			if err != nil {
				http.Error(w, err.Error(), 400)
			}
			pets = append(pets, &models.Pet{
				ID: ID,
				Name: name,
				Kind: kind,
			})
		}
		err = json.NewEncoder(w).Encode(pets)
		if err != nil {
			log.Print(err)
			http.Error(w, "Error while json-encoding", 400)
		}
	}
}

func generateNames(amount string) ([]string, error) {
	randomNameAPI += "/" + amount
	log.Println(randomNameAPI)
	resp, err := http.Get(randomNameAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var names []string
	log.Println(string(body))
	err = json.Unmarshal(body, &names)
	log.Print(names)
	if err != nil {
		return nil, err
	}
	return names, nil
}

func chooseYourDestiny() string {
	kinds := []string{"cat", "dog"}
	rand.Seed(time.Now().UnixNano())
	choosen := kinds[rand.Intn(2)]
	return choosen
}

func (env *Env) ListPetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pets, err := env.Store.ListPets()
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		json.NewEncoder(w).Encode(pets)
		if err != nil {
			log.Print(err)
			http.Error(w, "Error while json-encoding", 400)
		}
	}
}