package usecase

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"rincon-orlando/go-bootcamp/model"
)

type repo interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (*model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
}

// UseCase - Definition of a usecase layer, combining a repo and CSV filename
type UseCase struct {
	repo        repo
	csvFileName string
}

// Great help
// https://golangcode.com/how-to-read-a-csv-file-into-a-struct/

// New - UseCase factory
func New(repo repo, csvFilename string) (*UseCase, error) {
	// Open CSV file
	f, err := os.Open(csvFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read file into a variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	v := make([]model.Pokemon, 0, len(lines))

	// Build pokemons slice out of the file lines
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, errors.New("Error converting " + line[0] + " to int")
		}
		pokemon := model.Pokemon{
			ID:   id,
			Name: line[1],
		}
		v = append(v, pokemon)
	}

	// Build a new empty DB
	newUseCase := &UseCase{repo, csvFilename}
	// Then initialize the new DB with this particular set of Pokemons
	newUseCase.repo.SetPokemons(v)

	return newUseCase, nil
}

// GetAllPokemons - Returns a slice of all Pokemons available in this repository
func (uc UseCase) GetAllPokemons() []model.Pokemon {
	return uc.repo.GetAllPokemons()
}

// GetPokemonById - Returns a pokemon given its id
func (uc UseCase) GetPokemonById(id int) (*model.Pokemon, error) {
	return uc.repo.GetPokemonById(id)
}

// SetPokemons - Updates the internal repository with a new set of Pokemons
// Pointer as receiver so internal db can be modified
func (uc *UseCase) SetPokemons(pokemons []model.Pokemon) {
	uc.repo.SetPokemons(pokemons)
	// Once internal data is updated, persist it into the csv file
	uc.persist(pokemons)
}

func (uc UseCase) persist(pokemons []model.Pokemon) error {
	file, err := os.Create(uc.csvFileName)

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
