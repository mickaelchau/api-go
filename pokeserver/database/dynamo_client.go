package dynamo

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "pokedex"

type Pokemon struct {
	Name      string `json:"name"`
	Evolution int    `json:"evolution"`
	Poketype  string `json:"poketype"`
}

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

func DynamoGetResources() ([]Pokemon, error) {
	svc := InitSession()
	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Got error calling Scan: %s", err)
		return nil, err
	}
	var db_pokemon []Pokemon
	var pokemon Pokemon
	var name string
	var evolution int
	var poketype string
	for _, poke := range out.Items {
		name = *poke["name"].S //get the attribute "name" in map and get the str associated with it
		evolution, err = strconv.Atoi(*poke["evolution"].N)
		if err != nil {
			log.Printf("Got error calling Atoi: %s", err)
			return nil, err
		}
		poketype = *poke["poketype"].S
		pokemon.Name = name
		pokemon.Evolution = evolution
		pokemon.Poketype = poketype
		db_pokemon = append(db_pokemon, pokemon)
	}
	return db_pokemon, nil
}

func DynamoGetResource(pokeName string) (Pokemon, error) {
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
		log.Printf("Got error calling GetItem: %s", err)
		return pokemon, err
	}
	if result.Item == nil {
		return pokemon, nil
	}
	pokemon.Name = *result.Item["name"].S
	pokemon.Poketype = *result.Item["poketype"].S
	pokemon.Evolution, err = strconv.Atoi(*result.Item["evolution"].N)
	if err != nil {
		log.Printf("Got error calling Atoi: %s", err)
		return pokemon, err
	}
	return pokemon, nil
}

func DynamoDelete(name string) (bool, error) {
	svc := InitSession()
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := svc.DeleteItem(input)
	if err != nil {
		//status code
		log.Printf("Got error calling DeleteItem: %s", err)
		return false, err
	}
	return true, nil
}

func DynamoAdd(pokemon Pokemon) (bool, error) {
	svc := InitSession()
	serializedPokemon, err := dynamodbattribute.MarshalMap(pokemon) //map key: value
	if err != nil {
		log.Printf("Error: Marsahalling Pokemon failed")
		return false, err
	}

	input := &dynamodb.PutItemInput{
		Item:      serializedPokemon,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Got error calling PutItem: %s", err)
		return false, err
	}
	return true, nil
}
