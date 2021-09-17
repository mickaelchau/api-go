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

type AddResponse struct {
	Added   bool   `json:"added"`
	Pokemon string `json:"pokemon"`
}

type DeleteResponse struct {
	Deleted bool   `json:"deleted"`
	Pokemon string `json:"pokemon"`
}

func handleHome(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello, this application is an awesome Pokedex !")
}

func getAllPokemons(response http.ResponseWriter, request *http.Request) {
	var db_pokedex []dynamo.Pokemon = dynamo.DynamoGetPokedex()
	json.NewEncoder(response).Encode(db_pokedex)
}

func getSinglePokemon(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	dbResponse := dynamo.DynamoGetPokemon(name)
	json.NewEncoder(response).Encode(dbResponse)
}

func addToPokedex(response http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var pokemon dynamo.Pokemon
	json.Unmarshal(reqBody, &pokemon) //transform the JSON reqBody to a pokemon

	var dbResponse AddResponse
	dbResponse.Added = dynamo.DynamoAdd(pokemon)
	dbResponse.Pokemon = pokemon.Name
	json.NewEncoder(response).Encode(dbResponse)
}

func removeFromPokedex(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	var dbResponse DeleteResponse
	dbResponse.Deleted = dynamo.DynamoDelete(name)
	dbResponse.Pokemon = name
	json.NewEncoder(response).Encode(dbResponse)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/pokemons", getAllPokemons)
	router.HandleFunc("/pokemon/{name}", removeFromPokedex).Methods("DELETE")
	router.HandleFunc("/pokemon/{name}", getSinglePokemon)
	router.HandleFunc("/pokemon", addToPokedex).Methods("POST", "PUT")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	handleRequests()
}
