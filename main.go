package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mickael/api-server/dynamo"
	"net/http"

	"github.com/gorilla/mux"
)

func handleHome(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello, this application is an awesome Pokedex !")
}

func returnAllPokemons(response http.ResponseWriter, request *http.Request) {
	var db_pokedex []dynamo.Pokemon = dynamo.DynamoGetPokedex()
	json.NewEncoder(response).Encode(db_pokedex)
}

func returnSinglePokemon(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	dbResponse := dynamo.DynamoGetPokemon(name)
	json.NewEncoder(response).Encode(dbResponse)

}

func addToPokedex(response http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var pokemon dynamo.Pokemon
	json.Unmarshal(reqBody, &pokemon) //transform the JSON reqBody to a pokemon

	var dbResponse dynamo.AddResponse
	dbResponse.Added = dynamo.DynamoAdd(pokemon)
	dbResponse.Pokemon = pokemon.Name
	json.NewEncoder(response).Encode(dbResponse)
}

func removeFromPokedex(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	var dbResponse dynamo.DeleteResponse
	dbResponse.Deleted = dynamo.DynamoDelete(name)
	dbResponse.Pokemon = name
	json.NewEncoder(response).Encode(dbResponse)
}

/*
WITH DYNAMODB, ADD AND UPDATE ARE THE SAME
func updatePokedex(response http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var pokemon dynamo.Pokemon
	json.Unmarshal(reqBody, &pokemon)
	for _, db_pokemon := range Pokedex {
		if db_pokemon.Name == pokemon.Name {
			db_pokemon.Name = pokemon.Name
			db_pokemon.Evolution = pokemon.Evolution
			db_pokemon.Type = pokemon.Type
		}
	}
	json.NewEncoder(response).Encode(Pokedex)
}
*/

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/pokemons", returnAllPokemons)
	router.HandleFunc("/pokemon/{name}", removeFromPokedex).Methods("DELETE")
	router.HandleFunc("/pokemon/{name}", returnSinglePokemon)
	router.HandleFunc("/pokemon", addToPokedex).Methods("POST", "PUT")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	handleRequests()
}
