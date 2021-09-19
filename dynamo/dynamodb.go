package dynamo

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Pokemon struct {
	Name      string `json:"name"`
	Evolution int    `json:"evolution"`
	Poketype  string `json:"poketype"` //poketype is mapped to poketype in JSON
}

const tableName = "pokedex"

// USE OF THIS FUNC IS PROBABLY A BAD PRACTICE
func InitSession() *dynamodb.DynamoDB {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	return dynamodb.New(sess)
}

func DynamoGetPokemon(pokeName string) Pokemon {
	svc := InitSession()
	var pokemon Pokemon
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(pokeName),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	if result.Item == nil {
		return pokemon
	}
	pokemon.Name = *result.Item["name"].S
	pokemon.Poketype = *result.Item["poketype"].S
	pokemon.Evolution, err = strconv.Atoi(*result.Item["evolution"].N)
	if err != nil {
		log.Fatalf("Got error calling Atoi: %s", err)
	}
	return pokemon
}

func DynamoAdd(pokemon Pokemon) bool {
	svc := InitSession()
	serializedPokemon, err := dynamodbattribute.MarshalMap(pokemon) //map key: value
	if err != nil {
		log.Fatalf("Error: Marsahalling Pokemon failed")
		return false
	}

	input := &dynamodb.PutItemInput{
		Item:      serializedPokemon,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return false
	}
	return true
}

func DynamoDelete(pokeName string) bool {
	svc := InitSession()
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(pokeName),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Fatalf("Got error calling DeleteItem: %s", err)
		return false
	}
	return true
}

func DynamoGetPokedex() []Pokemon {
	svc := InitSession()
	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Fatalf("Got error calling Scan: %s", err)
	}
	var db_pokemon []Pokemon
	var pokemon Pokemon
	var name string
	var evolution int
	var poketype string
	for _, poke := range out.Items {
		name = *poke["name"].S //get the attribute "name" in map and get the str associated with it
		evolution, _ = strconv.Atoi(*poke["evolution"].N)
		poketype = *poke["poketype"].S
		pokemon.Name = name
		pokemon.Evolution = evolution
		pokemon.Poketype = poketype
		db_pokemon = append(db_pokemon, pokemon)
	}
	return db_pokemon
}
