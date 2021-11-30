package controller

import (
	"net/http"
	"rincon-orlando/go-bootcamp/model"
	"rincon-orlando/go-bootcamp/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

const pokemonApiUrl = "https://pokeapi.co/api/v2/pokemon/"

type repo interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
}

type Controller struct {
	repo repo
}

// Factory method
func NewController(repo repo) Controller {
	return Controller{repo}
}

// GetAllPokemons - service that returns all pokemons handled by the underlying repository
func (c Controller) GetAllPokemons(ctx *gin.Context) {
	data := c.repo.GetAllPokemons()
	ctx.IndentedJSON(http.StatusOK, data)
}

// GetPokemonById - service that returns a particular pokemon if it is present in the repository
func (c Controller) GetPokemonById(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to convert id " + id + " to int"})
		return
	}

	pokemon, err := c.repo.GetPokemonById(idInt)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, pokemon)
}

// GetPokemonById - service that fetchs a pokemon list from external api and persists it into the underlying repository
func (c Controller) FetchPokemonsFromApi(ctx *gin.Context) {
	data, err := util.FetchPokemonsFromApi(pokemonApiUrl)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Update repository underlying info
	c.repo.SetPokemons(data)

	ctx.IndentedJSON(http.StatusOK, c.repo.GetAllPokemons())
}
