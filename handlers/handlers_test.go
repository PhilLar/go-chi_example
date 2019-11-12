package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhilLar/go-chi_example/handlers"
	"github.com/PhilLar/go-chi_example/models"

	"github.com/gojuno/minimock/v3"
)

func TestListPetsHandler(t *testing.T) {
	t.Run("returns StatusOK", func(t *testing.T) {
		mc := minimock.NewController(t)
		defer mc.Finish()
		mockPetStore := handlers.NewPetStoreMock(mc)
		pets := []*models.Pet{
			&models.Pet{
				ID:   1,
				Name: "Barsik",
				Kind: "Cat",
				Age:  3,
			},
			&models.Pet{
				ID:   1,
				Name: "Jack",
				Kind: "Dog",
				Age:  10,
			},
			&models.Pet{
				ID:   1,
				Name: "Marsik",
				Kind: "Cat",
				Age:  7,
			},
		}

		mockPetStore.ListPetsMock.Return(pets, nil)

		env := &handlers.Env{Store: mockPetStore}
		req, err := http.NewRequest("GET", "/pets/get/all", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(env.ListPetsHandler())
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		expected := `[{"id":1,"first_name":"Krish","last_name":"Bhanushali","email_address":"krishsb@g.com","phone_number":"0987654321"},{"id":2,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"},{"id":6,"first_name":"FirstNameSample","last_name":"LastNameSample","email_address":"lr@gmail.com","phone_number":"1111111111"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}

	})
}
