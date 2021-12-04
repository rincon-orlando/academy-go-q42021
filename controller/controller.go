package controller

import (
	"net/http"
	"strconv"

	"rincon-orlando/go-bootcamp/model"

	"github.com/gin-gonic/gin"
)

type usecase interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (*model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
	FetchPokemonsFromApi() ([]model.Pokemon, error)
}

// Controller - Handler to communicate between endpoints and the usecase
type Controller struct {
	uc usecase
}

// New - Controller Factory
func New(uc usecase) Controller {
	return Controller{uc}
}

// GetAllPokemons - handler that returns all pokemons in the underlying repository
func (c Controller) GetAllPokemons(ctx *gin.Context) {
	data := c.uc.GetAllPokemons()
	ctx.IndentedJSON(http.StatusOK, data)
}

// GetPokemonById - handler that returns a particular pokemon if it is present in the underlying repository
func (c Controller) GetPokemonById(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to convert id " + id + " to int"})
		return
	}

	pokemon, err := c.uc.GetPokemonById(idInt)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, pokemon)
}

// FetchPokemonsFromApi - handlers that returns a pokemon list from external API
func (c Controller) FetchPokemonsFromApi(ctx *gin.Context) {
	data, err := c.uc.FetchPokemonsFromApi()
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Update repository underlying info
	c.uc.SetPokemons(data)

	ctx.IndentedJSON(http.StatusOK, c.uc.GetAllPokemons())
}
