package main

// Great help
// https://golangcode.com/how-to-read-a-csv-file-into-a-struct/

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const csvFilename = "pokemons.csv"

type Pokemon struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DB struct {
	data map[int]Pokemon
}

func main() {
	pokemonDB := new(DB)
	pokemonDB.initFromCsv(csvFilename)

	// Loop through pokemons and print 'em all
	for _, pokemon := range pokemonDB.data {
		fmt.Printf("ID: %v, Name %v\n", pokemon.Id, pokemon.Name)
	}

	// Configure router
	router := gin.Default()
	router.GET("/pokemons", pokemonDB.getAllPokemons)
	router.GET("/pokemons/:id", pokemonDB.getPokemonById)

	// Start server
	router.Run("localhost:8082")
}

// TODO: I do not like mixing DB and GIN context
// How do I make one method for getAllPokemons from DB
// and another one to just format as JSON?
// i.e. how do I share context between these two objects without
// making them dependent of each other?
func (db *DB) getAllPokemons(c *gin.Context) {
	// Convert DB map to array so we do not give unnecessary keys as output
	v := make([]Pokemon, 0, len(db.data))
	for _, value := range db.data {
		v = append(v, value)
	}
	c.IndentedJSON(http.StatusOK, v)
}

func (db *DB) getPokemonById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to convert id " + id + " to int"})
		return
	}

	if val, ok := db.data[idInt]; ok {
		c.IndentedJSON(http.StatusOK, val)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "pokemon with id " + id + " not found"})
}

func (db *DB) initFromCsv(filename string) (*DB, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return db, err
	}
	defer f.Close()

	// Read file into a variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return db, err
	}

	// Build pokemons out of the info
	pokemons := make(map[int]Pokemon)
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return db, errors.New("Error converting " + line[0] + " to int")
		}
		pokemon := Pokemon{
			Id:   id,
			Name: line[1],
		}
		pokemons[id] = pokemon
	}
	db.data = pokemons

	return db, nil
}
