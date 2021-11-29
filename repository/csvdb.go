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

// Great help
// https://golangcode.com/how-to-read-a-csv-file-into-a-struct/

// Factory method
func NewCSVDB(csvFilename string) (*CSVDB, error) {
	// Open CSV file
	f, err := os.Open(csvFilename)
	if err != nil {
		return &CSVDB{}, err
	}
	defer f.Close()

	// Read file into a variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return &CSVDB{}, err
	}

	v := make([]model.Pokemon, 0, len(lines))

	// Build pokemons slice out of the file lines
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return &CSVDB{}, errors.New("Error converting " + line[0] + " to int")
		}
		pokemon := model.Pokemon{
			ID:   id,
			Name: line[1],
		}
		v = append(v, pokemon)
	}

	// Build a new empty DB
	result := &CSVDB{DB{}, csvFilename}
	// Then initialize the new DB with this particular set of Pokemons
	result.db.setPokemons(v)

	return result, nil
}

// GetAllPokemons - Returns a slice of all Pokemons available in this repository
func (csvdb CSVDB) GetAllPokemons() []model.Pokemon {
	return csvdb.db.getAllPokemons()
}

// GetPokemonById - Returns a pokemon given its id
func (csvdb CSVDB) GetPokemonById(id int) (model.Pokemon, error) {
	return csvdb.db.getPokemonById(id)
}

// SetPokemons - Updates the internal repository with a new set of Pokemons
// Pointer as receiver so internal db can be modified
func (csvdb *CSVDB) SetPokemons(pokemons []model.Pokemon) {
	csvdb.db.setPokemons(pokemons)
	// Once internal data is updated, persist it into the csv file
	csvdb.persist()
}

func (csvdb CSVDB) persist() error {
	file, err := os.Create(csvdb.csvFileName)

	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range csvdb.db.data {
		line := []string{strconv.FormatInt(int64(value.ID), 10), value.Name}
		err := writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}
