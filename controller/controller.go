package controller

import (
	"net/http"
	"strconv"

	"rincon-orlando/go-bootcamp/model"
	"rincon-orlando/go-bootcamp/util/enum"

	"github.com/gin-gonic/gin"
)

type usecase interface {
	GetAllPokemons() []model.Pokemon
	GetPokemonById(id int) (*model.Pokemon, error)
	SetPokemons(pokemons []model.Pokemon)
	FetchPokemonsFromApi() ([]model.Pokemon, error)
	FilterPokemonsConcurrently(enum.OddEven, int, int, int) []model.Pokemon
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
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Failed to convert id " + id + " to int"})
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

// FilterPokemonsConcurrently - handler to return a list of odd/even pokemons processed concurrently
func (c Controller) FilterPokemonsConcurrently(ctx *gin.Context) {
	typeArg := ctx.Query("type")
	// Only support "odd" or "even". Parsing will know
	oddEven, err := enum.ParseOddEven(typeArg)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "'type' param error. " + err.Error()})
		return
	}

	// Amount of valid items you need to display as a response
	// TODO: Take the default value from the env
	items := ctx.DefaultQuery("items", "5")
	itemsInt, err := strconv.Atoi(items)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "'items' param error. Cannot convert " + items + " to int"})
		return
	}

	// Amount of valid items the worker should append to the response
	// TODO: Take the default value from the env
	ipw := ctx.DefaultQuery("items_per_workers", "10")
	ipwInt, err := strconv.Atoi(ipw)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "'items_per_workers' param error. Cannot convert " + ipw + " to int"})
		return
	}

	// As of now this is not a param.
	// TODO: Take the default value from the env
	numWorkers := 2

	ctx.IndentedJSON(http.StatusOK, c.uc.FilterPokemonsConcurrently(oddEven, numWorkers, itemsInt, ipwInt))
}
