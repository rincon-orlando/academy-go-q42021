package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"rincon-orlando/go-bootcamp/model"
)

// Useful doc
// https://tutorialedge.net/golang/consuming-restful-api-with-go/

// Response struct to map the entire Pokemon API response
type apiResponse struct {
	Results []apiPokemon `json:"results"`
}

// Pokemon struct to map only Pokemon entries from the respose (dismiss everything else)
type apiPokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Service - Definition of a service
type Service struct {
	url string
}

// New - Service factory
func New(url string) Service {
	return Service{url}
}

// FetchPokemonsFromApi - Utility method to try fetch Pokemons from a particular url
func (s Service) FetchPokemonsFromApi() ([]model.Pokemon, error) {
	response, err := http.Get(s.url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseObject apiResponse
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return nil, err
	}

	return mapPokemons(responseObject.Results), nil
}

// Utility method to extract the Pokemon ID out of the Url
func match(s string) (int, error) {
	i := strings.LastIndex(s, "/")
	if i >= 0 {
		j := strings.LastIndex(s[:i-1], "/")
		if j >= 0 {
			return strconv.Atoi(s[j+1 : i])
		}
	}

	return -1, nil
}

// Map API Pokemons into our own model Pokemons
func mapPokemons(pokemons []apiPokemon) []model.Pokemon {
	var v []model.Pokemon

	for _, pokemon := range pokemons {
		id, err := match(pokemon.Url)
		// This is an error we can tolerate, swallow it
		if err != nil {
			fmt.Printf("Error extracting Pokemon ID from url: %s", err)
			continue
		}

		v = append(v, model.Pokemon{ID: id, Name: pokemon.Name})
	}

	return v
}
