package service

import (
	"encoding/csv"
	"errors"
	"os"
	"rincon-orlando/go-bootcamp/model"
	"strconv"
)

type repo interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
}

type Service struct {
	repo        repo
	csvFileName string
}

// Great help
// https://golangcode.com/how-to-read-a-csv-file-into-a-struct/

// Factory method
func NewService(repo repo, csvFilename string) (*Service, error) {
	// Open CSV file
	f, err := os.Open(csvFilename)
	if err != nil {
		return &Service{}, err
	}
	defer f.Close()

	// Read file into a variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return &Service{}, err
	}

	v := make([]model.Pokemon, 0, len(lines))

	// Build pokemons slice out of the file lines
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return &Service{}, errors.New("Error converting " + line[0] + " to int")
		}
		pokemon := model.Pokemon{
			ID:   id,
			Name: line[1],
		}
		v = append(v, pokemon)
	}

	// Build a new empty DB
	newService := &Service{repo, csvFilename}
	// Then initialize the new DB with this particular set of Pokemons
	newService.repo.SetPokemons(v)

	return newService, nil
}

// GetAllPokemons - Returns a slice of all Pokemons available in this repository
func (s Service) GetAllPokemons() []model.Pokemon {
	return s.repo.GetAllPokemons()
}

// GetPokemonById - Returns a pokemon given its id
func (csvdb Service) GetPokemonById(id int) (model.Pokemon, error) {
	return csvdb.repo.GetPokemonById(id)
}

// SetPokemons - Updates the internal repository with a new set of Pokemons
// Pointer as receiver so internal db can be modified
func (s *Service) SetPokemons(pokemons []model.Pokemon) {
	s.repo.SetPokemons(pokemons)
	// Once internal data is updated, persist it into the csv file
	s.persist(pokemons)
}

func (s Service) persist(pokemons []model.Pokemon) error {
	file, err := os.Create(s.csvFileName)

	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range pokemons {
		line := []string{strconv.FormatInt(int64(value.ID), 10), value.Name}
		err := writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}
