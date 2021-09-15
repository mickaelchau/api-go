package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Pokemon struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Evolution int    `json:"Evolution"`
	Element   string `json:"Element"`
}

var Pokedex = []Pokemon{
	{Id: "1", Name: "Salamèche", Evolution: 1, Element: "Feu"},
	{Id: "2", Name: "Reptincel", Evolution: 2, Element: "Feu"},
	{Id: "3", Name: "Dracaufeu", Evolution: 3, Element: "Feu"},
	{Id: "4", Name: "Carapuce", Evolution: 1, Element: "Eau"},
	{Id: "5", Name: "Carabaffe", Evolution: 2, Element: "Eau"},
	{Id: "6", Name: "Tortank", Evolution: 3, Element: "Eau"},
	{Id: "7", Name: "Bulbizare", Evolution: 1, Element: "Plante"},
	{Id: "8", Name: "Herbizare", Evolution: 2, Element: "Plante"},
	{Id: "9", Name: "Florizare", Evolution: 3, Element: "Plante"},
}

func handleHome(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hello, this application is an awesome Pokedex !")
}

func returnAllPokemons(response http.ResponseWriter, request *http.Request) {
	fmt.Println("That's all the Pokemons!")
	json.NewEncoder(response).Encode(Pokedex)
}

func returnSinglePokemon(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["id"]

	for _, pokemon := range Pokedex {
		if pokemon.Id == key {
			json.NewEncoder(response).Encode(pokemon)
		}
	}
}

func addToPokedex(response http.ResponseWriter, request *http.Request) {
	/*if response.Header.Values("Content-Type")[0] != "application/json" {
		content.WriteHeader(http.StatusBadGateway)
	} */
	reqBody, _ := ioutil.ReadAll(request.Body)
	var pokemon Pokemon
	json.Unmarshal(reqBody, &pokemon) //transform the JSON reqBody to a pokemon

	Pokedex = append(Pokedex, pokemon)
	pokemon.Id = fmt.Sprintf("%d", len(Pokedex))
	json.NewEncoder(response).Encode(pokemon) //Response the object pokemon (does not make sense, we should send an object ID)
}

func removeFromPokedex(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]

	for index, pokemon := range Pokedex {
		if pokemon.Id == id {
			Pokedex = append(Pokedex[:index], Pokedex[index+1:]...)
		}
	}
	json.NewEncoder(response).Encode(Pokedex)
}
func updatePokedex(response http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var pokemon Pokemon
	json.Unmarshal(reqBody, &pokemon)
	for _, db_pokemon := range Pokedex {
		if db_pokemon.Id == pokemon.Id {
			db_pokemon.Name = pokemon.Name
			db_pokemon.Evolution = pokemon.Evolution
			db_pokemon.Element = pokemon.Element
			log.Printf("%s", pokemon.Name)
		}
	}
	fmt.Println("testst")
	json.NewEncoder(response).Encode(Pokedex)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handleHome)
	router.HandleFunc("/pokemons", returnAllPokemons)
	router.HandleFunc("/pokemon/{id}", removeFromPokedex).Methods("DELETE")
	router.HandleFunc("/pokemon/{id}", returnSinglePokemon)
	router.HandleFunc("/pokemon", addToPokedex).Methods("POST")
	router.HandleFunc("/pokemon", updatePokedex).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	handleRequests()
}
