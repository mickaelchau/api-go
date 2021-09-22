package pokeservice

import (
	"log"

	dynamo "gitlab.com/tcmlabs/api-webserver/pokeserver/database"
)

func GetAllPokemons() ([]dynamo.Pokemon, error) {
	var pokedex []dynamo.Pokemon
	var err error
	pokedex, err = dynamo.DynamoGetResources()
	if err != nil {
		log.Printf("Got error calling DynamoGetResources: %s", err)
		return []dynamo.Pokemon{}, err
	}
	return pokedex, nil //err != nil => err ==nil
}

func GetSinglePokemon(name string) (dynamo.Pokemon, error) {
	var pokemon dynamo.Pokemon
	var err error
	pokemon, err = dynamo.DynamoGetResource(name)
	if err != nil {
		log.Printf("Got error calling DynamoGetResource: %s", err)
		return dynamo.Pokemon{}, err
	}
	return pokemon, nil
}

func AddToPokedex(pokemon dynamo.Pokemon) error {
	return dynamo.DynamoAdd(pokemon)
}

func RemoveFromPokedex(name string) error {
	return dynamo.DynamoDelete(name)
}
