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

type InvalidResponse struct {
	Error string `json:"error"`
}

const invalidPathOrUrl = "Invalid URL."
const invalidHeaders = "'Content-Type' or/and 'Accept' Header(s) are not set to 'application/json'."
const invalidPokemonName = "The pokemon you try to get is not in the database."
const invalidBody = "The pokemon you try to add must have a name (string), a poketype (string) and an evolution (int)."

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
	foundPokemon := dynamo.DynamoGetPokemon(name)
	if foundPokemon.Name == "" {
		handleError(response, invalidPokemonName, http.StatusNotFound)
	} else {
		json.NewEncoder(response).Encode(foundPokemon)
	}
}

func addToPokedex(response http.ResponseWriter, request *http.Request) {
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatalf("Got error calling ReadAll: %s", err)
	}
	var pokemon dynamo.Pokemon
	json.Unmarshal(reqBody, &pokemon) //transform the JSON reqBody to a pokemon
	if pokemon.Poketype == "" || pokemon.Name == "" || pokemon.Evolution == 0 {
		handleError(response, invalidBody, http.StatusBadRequest)
	} else {
		var dbResponse AddResponse
		dbResponse.Added = dynamo.DynamoAdd(pokemon)
		dbResponse.Pokemon = pokemon.Name
		json.NewEncoder(response).Encode(dbResponse)
	}
}

func removeFromPokedex(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	var dbResponse DeleteResponse
	dbResponse.Deleted = dynamo.DynamoDelete(name)
	dbResponse.Pokemon = name
	json.NewEncoder(response).Encode(dbResponse)
}

func handleError(response http.ResponseWriter, errorDescription string, errorCode int) {
	http.Error(response, errorDescription, errorCode)
}

func handleInvalidURL(response http.ResponseWriter, request *http.Request) {
	handleError(response, invalidPathOrUrl, http.StatusForbidden)
}

func headerCheckMiddleWare(finalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		url := request.RequestURI
		content := request.Header.Get("Content-Type")
		accept := request.Header.Get("Accept")
		if url == "/" || (content == "application/json" && accept == "application/json") {
			finalHandler.ServeHTTP(response, request)
		} else {
			handleError(response, invalidHeaders, http.StatusNotAcceptable)
		}
	})
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(headerCheckMiddleWare) //will be executed ONLY if a routre is matched
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/pokemons", getAllPokemons).Methods("GET")
	router.HandleFunc("/pokemon/{name}", removeFromPokedex).Methods("DELETE")
	router.HandleFunc("/pokemon/{name}", getSinglePokemon).Methods("GET")
	router.HandleFunc("/pokemon", addToPokedex).Methods("POST", "PUT")
	router.NotFoundHandler = http.HandlerFunc(handleInvalidURL)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	handleRequests()
}
