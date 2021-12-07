package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"rincon-orlando/go-bootcamp/model"

	"github.com/stretchr/testify/assert"
)

const validFakeResponse string = `
{
	"count": 1118,
	"next": "https://pokeapi.co/api/v2/pokemon/?offset=20&limit=20",
	"previous": null,
	"results": [
	{
	"name": "bulbasaur",
	"url": "https://pokeapi.co/api/v2/pokemon/1/"
	},
	{
	"name": "ivysaur",
	"url": "https://pokeapi.co/api/v2/pokemon/2/"
	}
	]
	}
`

const invalidFakeResponse string = `
{
	"count": 1118,
	"next": "https://pokeapi.co/api/v2/pokemon/?offset=20&limit=20",
	"previous": null,
	"results": [
	{
	"name": "bulbasaur",
	"url": "https://pokeapi.co/api/v2/pokemon/1/"
	},
	{
	"name": "ivysaur",
	"url": "https://pokeapi.co/api/v2/pokemon/2/"
	}
	]
`

// TestService_FetchPokemonsFromApi - Test client to get pokemons from external API
func TestService_FetchPokemonsFromApi(t *testing.T) {
	testCases := []struct {
		name             string
		serverResponse   string
		expectedPokemons []model.Pokemon
		hasError         bool
		error            error
	}{
		{
			name:           "fetch pokemon OK",
			serverResponse: validFakeResponse,
			expectedPokemons: []model.Pokemon{
				{
					ID: 1, Name: "bulbasaur",
				},
				{
					ID: 2, Name: "ivysaur",
				},
			},
			hasError: false,
			error:    nil,
		},
		{
			name:             "fetch pokemon wrong response",
			serverResponse:   invalidFakeResponse,
			expectedPokemons: []model.Pokemon{},
			hasError:         true,
			error:            errors.New("unexpected end of JSON input"),
		},
	}

	for _, tc := range testCases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, tc.serverResponse)
		}))

		defer server.Close()

		service := New(server.URL)

		pokemons, err := service.FetchPokemonsFromApi()
		if tc.hasError {
			assert.EqualError(t, err, tc.error.Error())
		} else {
			assert.EqualValues(t, tc.expectedPokemons, pokemons)
		}
	}

}
