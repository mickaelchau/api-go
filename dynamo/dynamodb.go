package dynamo

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Pokemon struct {
	Name      string `json:"name"`
	Evolution int    `json:"evolution"`
	Poketype  string `json:"poketype"` //Pokepokepoketype is mapped to pokepoketype in JSON
}

type AddResponse struct {
	Added   bool   `json:"added"`
	Pokemon string `json:"pokemon"`
}

type DeleteResponse struct {
	Deleted bool   `json:"deleted"`
	Pokemon string `json:"pokemon"`
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

func InitEc2() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ec2.New(sess)

	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String("ami-072056ff9d3689e7b"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}
	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)
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
		log.Fatalf("Got error calling DeleteItem: %s", err)
	}
	if result.Item == nil {
		return pokemon
	}
	pokemon.Name = *result.Item["name"].S
	pokemon.Poketype = *result.Item["poketype"].S
	pokemon.Evolution, _ = strconv.Atoi(*result.Item["evolution"].N)
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
