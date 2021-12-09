package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var evenPokemon = Pokemon{
	ID:   8,
	Name: "charmander",
}

var oddPokemon = Pokemon{
	ID:   9,
	Name: "squirtle",
}

// TestPokemon_IsEven - test the IsEven method from Pokemon
func TestPokemon_IsEven(t *testing.T) {
	testCases := []struct {
		name           string
		pokemon        Pokemon
		expectedResult bool
	}{
		{
			name:           "test even pokemon",
			pokemon:        evenPokemon,
			expectedResult: true,
		},
		{
			name:           "test odd pokemon",
			pokemon:        oddPokemon,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedResult, tc.pokemon.IsEven())
	}
}

// TestPokemon_String - test the String method from Pokemon
func TestPokemon_String(t *testing.T) {
	testCases := []struct {
		name           string
		pokemon        Pokemon
		expectedResult string
	}{
		{
			name:           "test pokemon to string",
			pokemon:        evenPokemon,
			expectedResult: "ID = 8, Name = charmander",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expectedResult, tc.pokemon.String())
	}
}
