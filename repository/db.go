package repository

import (
	"fmt"
	"rincon-orlando/go-bootcamp/model"
)

// DB - Definition of a Pokemon repository, simulating a Database
type DB struct {
	data map[int]model.Pokemon
}

// getAllPokemons - Returns a slice of all Pokemons available in this repository
func (db DB) getAllPokemons() []model.Pokemon {
	// Convert map to array so we do not give unnecessary keys as output
	v := make([]model.Pokemon, 0, len(db.data))
	for _, value := range db.data {
		v = append(v, value)
	}
	return v
}

// getPokemonById - Returns a pokemon given its id
func (db DB) getPokemonById(id int) (model.Pokemon, error) {
	if val, ok := db.data[id]; ok {
		return val, nil
	}

	return model.Pokemon{}, fmt.Errorf("pokemon with %d not found", id)
}

// setPokemons - Build a pokemon map out of the pokemon slice
// Pointer as receiver so internal db.data can be modified
func (db *DB) setPokemons(pokemons []model.Pokemon) {
	db.data = make(map[int]model.Pokemon)
	for _, pokemon := range pokemons {
		db.data[pokemon.ID] = pokemon
	}
}
