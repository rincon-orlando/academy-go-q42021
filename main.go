package main

// Great help
// https://golangcode.com/how-to-read-a-csv-file-into-a-struct/

import (
	"net/http"
	"rincon-orlando/go-bootcamp/model"
	"rincon-orlando/go-bootcamp/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

const csvFilename = "pokemons.csv"

type PokemonDB interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (model.Pokemon, error)
}

var pokemonDB PokemonDB

func main() {
	var err error
	pokemonDB, err = repository.NewCSVDB(csvFilename)
	if err != nil {
		// TODO: Handle error here
	}

	// Configure router
	router := gin.Default()
	router.GET("/pokemons", getAllPokemons)
	router.GET("/pokemons/:id", getPokemonById)

	// Start server
	router.Run("localhost:8082")
}

func getAllPokemons(c *gin.Context) {
	data := pokemonDB.GetAllPokemons()
	c.IndentedJSON(http.StatusOK, data)
}

func getPokemonById(c *gin.Context) {
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
