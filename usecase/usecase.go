package usecase

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"rincon-orlando/go-bootcamp/config"
	"rincon-orlando/go-bootcamp/model"
	"rincon-orlando/go-bootcamp/util/enum"
	"rincon-orlando/go-bootcamp/workerpool"
)

type repo interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (*model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
}

type service interface {
	FetchPokemonsFromApi() ([]model.Pokemon, error)
}

// UseCase - Definition of a usecase layer, combining a repo, a csv filename + service (external API client)
type UseCase struct {
	repo        repo
	csvFileName string
	service     service
}

// Great help
// https://golangcode.com/how-to-read-a-csv-file-into-a-struct/

// New - UseCase factory
func New(repo repo, csvFilename string, service service) (*UseCase, error) {
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
	newUseCase := &UseCase{repo, csvFilename, service}
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

// FetchPokemonsFromApi - Returns a slice of Pokemons from external API
func (uc UseCase) FetchPokemonsFromApi() ([]model.Pokemon, error) {
	return uc.service.FetchPokemonsFromApi()
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

type pokeTask struct {
	pokemon   model.Pokemon
	processor func(model.Pokemon) bool
}

// Run - Poketask method execution, ready to send results to input channel
func (pt pokeTask) Run(ch chan<- interface{}) bool {
	if pt.processor(pt.pokemon) {
		// Inject the pokemon to the output channel
		ch <- pt.pokemon
		return true
	}
	return false
}

// FilterPokemonsConcurrently - Configures and executes a worker pool to extract a set of pokemons based off the given criteria
func (uc UseCase) FilterPokemonsConcurrently(oe enum.OddEven, numWorkers int, items int, ipw int) []model.Pokemon {
	// First, get all pokemons to work
	allPokemons := uc.GetAllPokemons()

	fmt.Printf("Worker config: numWorkers %d, items = %d, items_per_worker = %d\n", numWorkers, items, ipw)

	// TODO: I may want to move the dependency injection somewhere else
	config := config.NewPoolConfig(oe, numWorkers, items, ipw)
	pool := workerpool.New(numWorkers, config)

	// Process a Pokemon, to verify if matches what we are looking for
	filterPokemon := func(pokemon model.Pokemon) bool {
		if oe == enum.Even {
			return pokemon.IsEven()
		}
		return !pokemon.IsEven()
	}

	// Task config
	var tasks []pokeTask
	for _, p := range allPokemons {
		tasks = append(tasks, pokeTask{
			pokemon:   p,
			processor: filterPokemon,
		})
	}

	// Scheduling work
	for _, task := range tasks {
		pool.ScheduleWork(task)
	}

	// Worker pool returns a generic interface element
	genericResponse := pool.Monitor()

	// So, turn those interfaces into pokemons
	response := make([]model.Pokemon, len(genericResponse))
	for i, v := range genericResponse {
		response[i] = v.(model.Pokemon)
	}

	pool.Close()

	return response
}
