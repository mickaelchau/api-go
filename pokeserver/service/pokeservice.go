package pokeservice

import (
	dynamo "gitlab.com/tcmlabs/api-webserver/pokeserver/database"
)

func GetAllPokemons() []dynamo.Pokemon {
	var pokedex []dynamo.Pokemon = dynamo.DynamoGetPokedex()
	return pokedex
}

func GetSinglePokemon(name string) dynamo.Pokemon {
	var pokemon dynamo.Pokemon = dynamo.DynamoGetPokemon(name)
	return pokemon
}

func AddToPokedex(pokemon dynamo.Pokemon) bool {
	return dynamo.DynamoAdd(pokemon)
}

func RemoveFromPokedex(name string) bool {
	return dynamo.DynamoDelete(name)
}
