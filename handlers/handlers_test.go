package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhilLar/go-chi_example/handlers"
	"github.com/PhilLar/go-chi_example/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"

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
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(env.ListPetsHandler())
		handler.ServeHTTP(rec, req)

		expected := []models.Pet{
			models.Pet{
				ID:   1,
				Name: "Barsik",
				Kind: "Cat",
				Age:  3,
			},
			models.Pet{
				ID:   1,
				Name: "Jack",
				Kind: "Dog",
				Age:  10,
			},
			models.Pet{
				ID:   1,
				Name: "Marsik",
				Kind: "Cat",
				Age:  7,
			},
		}
		var actual []models.Pet
		assert.Equal(t, http.StatusOK, rec.Code)
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
		assert.Equal(t, actual, expected)
	})
}
