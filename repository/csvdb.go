package repository

import (
	"encoding/csv"
	"errors"
	"os"
	"rincon-orlando/go-bootcamp/model"
	"strconv"
)

type CSVDB struct {
	db          DB
	csvFileName string
}

// Factory method
func NewCSVDB(csvFilename string) (CSVDB, error) {
	// Open CSV file
	f, err := os.Open(csvFilename)
	if err != nil {
		return CSVDB{}, err
	}
	defer f.Close()

	// Read file into a variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return CSVDB{}, err
	}

	// Build pokemons out of the info
	pokemons := make(map[int]model.Pokemon)
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return CSVDB{}, errors.New("Error converting " + line[0] + " to int")
		}
		pokemon := model.Pokemon{
			ID:   id,
			Name: line[1],
		}
		pokemons[id] = pokemon
	}

	return CSVDB{DB{pokemons}, csvFilename}, nil
}

// GetAllPokemons - Returns a slice of all Pokemons available in this repository
func (csvdb CSVDB) GetAllPokemons() []model.Pokemon {
	return csvdb.db.getAllPokemons()
}

// GetPokemon - Returns a pokemon given its id
func (csvdb CSVDB) GetPokemonById(id int) (model.Pokemon, error) {
	return csvdb.db.getPokemonById(id)
}

func (csvdb CSVDB) Persist() error {
	// TODO: Implement
	return nil
}
