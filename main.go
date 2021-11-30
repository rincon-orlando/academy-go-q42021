package main

import (
	"log"
	"rincon-orlando/go-bootcamp/controller"
	"rincon-orlando/go-bootcamp/repository"
	"rincon-orlando/go-bootcamp/service"

	"github.com/gin-gonic/gin"
)

const csvFilename = "pokemons.csv"

func main() {
	// Dependency injection
	db := repository.NewDB()
	service, err := service.NewService(&db, csvFilename)
	if err != nil {
		log.Fatal("Error starting up database" + err.Error())
	}
	controller := controller.NewController(service)

	// Configure router
	router := gin.Default()
	router.GET("/pokemons", controller.GetAllPokemons)
	router.GET("/pokemons/:id", controller.GetPokemonById)
	router.GET("/pokemons/fetch", controller.FetchPokemonsFromApi)

	// Start server
	router.Run("localhost:8082")
}
