package main

import (
	"rincon-orlando/go-bootcamp/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configure router
	router := gin.Default()
	router.GET("/pokemons", controller.GetAllPokemons)
	router.GET("/pokemons/:id", controller.GetPokemonById)
	router.GET("/pokemons/fetch", controller.FetchPokemonsFromApi)

	// Start server
	router.Run("localhost:8082")
}
