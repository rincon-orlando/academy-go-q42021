package repository

import (
	"errors"
	"testing"

	"rincon-orlando/go-bootcamp/model"

	"github.com/stretchr/testify/assert"
)

var pokemons = []model.Pokemon{
	{ID: 1, Name: "Onix"},
	{ID: 2, Name: "Mewtwo"},
	{ID: 5, Name: "Pikachu"},
	{ID: 11, Name: "Bulbasur"},
	{ID: 100, Name: "Snorlax"},
}

// TestDB_New - Test DB factory method
func TestDB_New(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "factory method",
		},
	}

	// There is really one test to execute
	db := New()
	for range testCases {
		assert.Zero(t, len(db.pokeMap))
	}
}

// TestDB_SetAllPokemons - Test setting all pokemons and verify the internal structure
func TestDB_SetAllPokemons(t *testing.T) {
	testCases := []struct {
		name            string
		expectedId      int
		expectedPokemon model.Pokemon
	}{
		{
			name:            "test Onix is set",
			expectedId:      1,
			expectedPokemon: model.Pokemon{ID: 1, Name: "Onix"},
		},
		{
			name:            "test Mewtwo is set",
			expectedId:      2,
			expectedPokemon: model.Pokemon{ID: 2, Name: "Mewtwo"},
		},
		{
			name:            "test Pikachu is set",
			expectedId:      5,
			expectedPokemon: model.Pokemon{ID: 5, Name: "Pikachu"},
		},
		{
			name:            "test Bulbasur is set",
			expectedId:      11,
			expectedPokemon: model.Pokemon{ID: 11, Name: "Bulbasur"},
		},
		{
			name:            "test Snorlax is set",
			expectedId:      100,
			expectedPokemon: model.Pokemon{ID: 100, Name: "Snorlax"},
		},
	}

	db := New()
	db.SetPokemons(pokemons)
	for _, tc := range testCases {
		assert.Contains(t, db.pokeMap, tc.expectedId)
		assert.EqualValues(t, tc.expectedPokemon, db.pokeMap[tc.expectedId])
	}
}

// TestDB_GetAllPokemons - Test obtaining all pokemons
func TestDB_GetAllPokemons(t *testing.T) {
	testCases := []struct {
		name           string
		givenPokemons  []model.Pokemon
		expectedLength int
	}{
		{
			name:           "get all pokemons",
			expectedLength: 5,
			givenPokemons:  pokemons,
		},
		{
			name:           "get all pokemons with less data",
			expectedLength: 3,
			givenPokemons:  pokemons[:3],
		},
	}

	db := New()
	for _, tc := range testCases {
		db.SetPokemons(tc.givenPokemons)
		assert.Equal(t, tc.expectedLength, len(db.GetAllPokemons()))
	}
}

// TestDB_GetPokemonById - Test DB get pokemons by id
func TestDB_GetPokemonById(t *testing.T) {
	testCases := []struct {
		name        string
		id          int
		pokemonName string
		hasError    bool
		error       error
	}{
		{
			name:        "get existing pokemon",
			id:          1,
			pokemonName: "Onix",
			hasError:    false,
			error:       nil,
		},
		{
			name:     "get unexisting pokemon",
			id:       3,
			hasError: true,
			error:    errors.New("pokemon with 3 not found"),
		},
	}

	db := New()
	db.SetPokemons(pokemons)
	for _, tc := range testCases {
		pokemon, err := db.GetPokemonById(tc.id)
		if tc.hasError {
			assert.EqualError(t, err, tc.error.Error())
		} else {
			assert.Equal(t, tc.id, pokemon.ID)
			assert.Equal(t, tc.pokemonName, pokemon.Name)
		}
	}
}
