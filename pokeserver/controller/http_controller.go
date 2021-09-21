package http_controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dynamo "gitlab.com/tcmlabs/api-webserver/pokeserver/database"
	service "gitlab.com/tcmlabs/api-webserver/pokeserver/service"
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
const applicationName = "Pokédex"
const internalServerError = "Oops ! Something wrong happens in the server-side"

func HandleHome(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello, this application is an awesome %s!", applicationName)
}

func handleGets(response http.ResponseWriter, request *http.Request) {
	var getsResult []dynamo.Pokemon
	var err error
	getsResult, err = service.GetAllPokemons()
	if getsResult == nil || err != nil {
		handleError(response, internalServerError, http.StatusInternalServerError)
	}
	json.NewEncoder(response).Encode(getsResult)
}

func handleGet(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	var pokemon dynamo.Pokemon
	var err error
	pokemon, err = service.GetSinglePokemon(name)
	if err != nil {
		handleError(response, internalServerError, http.StatusInternalServerError)
	} else if pokemon.Name == "" {
		handleError(response, invalidPokemonName, http.StatusNotFound)
	} else {
		json.NewEncoder(response).Encode(pokemon)
	}
}

func handleDelete(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]
	var deleteResponse DeleteResponse
	var err error
	deleteResponse.Deleted, err = service.RemoveFromPokedex(name)
	if err != nil {
		handleError(response, internalServerError, http.StatusInternalServerError)
	} else {
		deleteResponse.Pokemon = name
		json.NewEncoder(response).Encode(deleteResponse)
	}
}

func HandlePostAndPut(response http.ResponseWriter, request *http.Request) {
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Got error calling ReadAll: %s", err)
		handleError(response, internalServerError, http.StatusInternalServerError)
	}
	var pokebody dynamo.Pokemon
	json.Unmarshal(reqBody, &pokebody) //transform the JSON reqBody to a pokemon
	if pokebody.Name == "" || pokebody.Poketype == "" || pokebody.Evolution == 0 {
		handleError(response, invalidBody, http.StatusBadRequest)
	} else {
		var addResponse AddResponse
		var err error
		addResponse.Added, err = service.AddToPokedex(pokebody)
		if err != nil {
			handleError(response, invalidBody, http.StatusBadRequest)
		}
		addResponse.Pokemon = pokebody.Name
		json.NewEncoder(response).Encode(addResponse)
	}
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

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(headerCheckMiddleWare) //will be executed ONLY if a routre is matched
	router.HandleFunc("/", HandleHome)
	router.HandleFunc("/pokemons", handleGets).Methods("GET")
	router.HandleFunc("/pokemon/{name}", handleGet).Methods("GET")
	router.HandleFunc("/pokemon/{name}", handleDelete).Methods("DELETE")
	router.HandleFunc("/pokemon", HandlePostAndPut).Methods("POST", "PUT")
	router.NotFoundHandler = http.HandlerFunc(handleInvalidURL)
	log.Fatal(http.ListenAndServe(":8000", router))
}
