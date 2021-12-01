package main

import (
	"log"

	"rincon-orlando/go-bootcamp/config"
	"rincon-orlando/go-bootcamp/controller"
	"rincon-orlando/go-bootcamp/repository"
	"rincon-orlando/go-bootcamp/service"
	"rincon-orlando/go-bootcamp/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.New(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Dependency injection
	db := repository.New()
	usecase, err := usecase.New(&db, config.CSV_FILENAME)
	if err != nil {
		log.Fatal("Error starting up database" + err.Error())
	}
	service := service.New(config.POKEMON_API_URL)

	controller := controller.New(usecase, service)

	// Configure router
	router := gin.Default()
	router.GET("/pokemons", controller.GetAllPokemons)
	router.GET("/pokemons/:id", controller.GetPokemonById)
	router.GET("/pokemons/fetch", controller.FetchPokemonsFromApi)

	// Start server
	router.Run("localhost:8082")
}
