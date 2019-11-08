package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PhilLar/go-chi_example/models"
	"github.com/go-chi/chi"
)

var randomNameAPI string = "http://names.drycodes.com"

type PetStore interface {
	InsertPet(age int, name, kind string) (int, error)
	ListPets() ([]*models.Pet, error)
	FilterPets(kind string, firstLetter string, underage, overage int) ([]*models.Pet, error)
	RemoveAllPets() error
}

type Env struct {
	Store PetStore
}

func GetPetHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	w.Write([]byte(fmt.Sprintf("get pet with name: %s", name)))
}

func (env *Env) PutPetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kind := chi.URLParam(r, "kind")
		name := chi.URLParam(r, "name")
		age := chi.URLParam(r, "age")
		ageInt, err := strconv.Atoi(age)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), 400)
		}
		ID, err := env.Store.InsertPet(ageInt, name, kind)
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
		for i := 0; i < amount; i++ {
			name := names[i]
			kind := chooseYourDestiny()
			age := generateAge()
			ID, err := env.Store.InsertPet(age, name, kind)
			if err != nil {
				http.Error(w, err.Error(), 400)
			}
			pets = append(pets, &models.Pet{
				ID:   ID,
				Name: name,
				Kind: kind,
				Age:  age,
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

func generateAge() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 15
	return rand.Intn(max-min+1) + min
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

func (env *Env) FilterPetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		underage := chi.URLParam(r, "underage")
		overage := chi.URLParam(r, "overage")
		underageInt, err := strconv.Atoi(underage)
		overageInt, err := strconv.Atoi(overage)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		kind := chi.URLParam(r, "kind")
		firstLetter := chi.URLParam(r, "first_letter")
		pets, err := env.Store.FilterPets(kind, firstLetter, underageInt, overageInt)
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

func (env *Env) RemoveAllPetsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := env.Store.RemoveAllPets()
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
	}
}
