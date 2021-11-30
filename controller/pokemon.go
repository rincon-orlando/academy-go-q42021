package controller

import (
	"net/http"
	"rincon-orlando/go-bootcamp/model"
	"rincon-orlando/go-bootcamp/repository"
	"rincon-orlando/go-bootcamp/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

const csvFilename = "pokemons.csv"
const pokemonApiUrl = "https://pokeapi.co/api/v2/pokemon/"

type PokemonDB interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
}

var pokemonDB PokemonDB

// TODO: I am not convinced I will be able to unit test this
// initPokemonDB builds a new DB to be shared across different service handlers
func initPokemonDB() PokemonDB {
	// Init shared pokemonDB at startup
	db, err := repository.NewCSVDB(csvFilename)
	if err != nil {
		// TODO: Handle error here
	}
	return db
}

// GetAllPokemons - service that returns all pokemons handled by the underlying repository
func GetAllPokemons(c *gin.Context) {
	if pokemonDB == nil {
		pokemonDB = initPokemonDB()
	}

	data := pokemonDB.GetAllPokemons()
	c.IndentedJSON(http.StatusOK, data)
}

// GetPokemonById - service that returns a particular pokemon if it is present in the repository
func GetPokemonById(c *gin.Context) {
	if pokemonDB == nil {
		pokemonDB = initPokemonDB()
	}

	id := c.Param("id")
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to convert id " + id + " to int"})
		return
	}

	pokemon, err := pokemonDB.GetPokemonById(idInt)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, pokemon)
}

// GetPokemonById - service that fetchs a pokemon list from external api and persists it into the underlying repository
func FetchPokemonsFromApi(c *gin.Context) {
	if pokemonDB == nil {
		pokemonDB = initPokemonDB()
	}

	data, err := util.FetchPokemonsFromApi(pokemonApiUrl)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Update repository underlying info
	pokemonDB.SetPokemons(data)

	c.IndentedJSON(http.StatusOK, pokemonDB.GetAllPokemons())
}
