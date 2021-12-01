package repository

import (
	"fmt"

	"rincon-orlando/go-bootcamp/model"
)

// DB - Definition of a Pokemon repository, simulating a Database
type DB struct {
	pokeMap map[int]model.Pokemon
}

// New - DB factory
func New() DB {
	return DB{}
}

// GetAllPokemons - Returns a slice of all Pokemons available in this repository
func (db DB) GetAllPokemons() []model.Pokemon {
	// Convert map to array so we do not give unnecessary keys as output
	v := make([]model.Pokemon, 0, len(db.pokeMap))
	for _, value := range db.pokeMap {
		v = append(v, value)
	}
	return v
}

// GetPokemonById - Returns a pokemon given its id
func (db DB) GetPokemonById(id int) (*model.Pokemon, error) {
	if val, ok := db.pokeMap[id]; ok {
		return &val, nil
	}

	return nil, fmt.Errorf("pokemon with %d not found", id)
}

// SetPokemons - Build a pokemon map out of the pokemon slice
// Pointer as receiver so internal db.data can be modified
func (db *DB) SetPokemons(pokemons []model.Pokemon) {
	db.pokeMap = make(map[int]model.Pokemon)
	for _, pokemon := range pokemons {
		db.pokeMap[pokemon.ID] = pokemon
	}
}
