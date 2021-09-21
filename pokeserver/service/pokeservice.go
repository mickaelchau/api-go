package pokeservice

import (
	dynamo "gitlab.com/tcmlabs/api-webserver/pokeserver/database"
)

func GetAllPokemons() ([]dynamo.Pokemon, error) {
	var pokedex []dynamo.Pokemon
	var err error
	pokedex, err = dynamo.DynamoGetResources()
	if err != nil {
		return nil, err
	}
	return pokedex, err
}

func GetSinglePokemon(name string) (dynamo.Pokemon, error) {
	var pokemon dynamo.Pokemon
	var err error
	pokemon, err = dynamo.DynamoGetResource(name)
	return pokemon, err
}

func AddToPokedex(pokemon dynamo.Pokemon) (bool, error) {
	return dynamo.DynamoAdd(pokemon)
}

func RemoveFromPokedex(name string) (bool, error) {
	return dynamo.DynamoDelete(name)
}
