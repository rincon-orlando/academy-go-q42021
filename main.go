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
	// Clean architecture layer order:
	// router -> controller -> usecase -> service / repository
	db := repository.New()
	service := service.New(config.POKEMON_API_URL)
	usecase, err := usecase.New(&db, config.CSV_FILENAME, service)
	if err != nil {
		log.Fatal("Error starting up database" + err.Error())
	}

	controller := controller.New(usecase)

	// Configure router
	router := gin.Default()
	router.GET("/pokemons", controller.GetAllPokemons)
	router.GET("/pokemons/:id", controller.GetPokemonById)
	router.GET("/pokemons/fetch", controller.FetchPokemonsFromApi)

	// Start server
	router.Run("localhost:8082")
}
