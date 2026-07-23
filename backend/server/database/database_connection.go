package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client{ 
	//must use full relative path for .Load()
	err := godotenv.Load("/Users/kevindev/Desktop/movienightv2/backend/.env")
	//checking to see if there is an error
	if err != nil {
		log.Println("Warning: unable to find .env file!")
	}
	//if no error then grab environment variable
	MongoDB := os.Getenv("MONGODB_URI")
	//if environment variable is empty, log fatal error
	if MongoDB == "" {
		log.Fatal("MONGODB_URI not set!")
	}
	//if not empty, print to console and add to clientOptions function
	fmt.Println("MongoDB URI: ", MongoDB)
	//create clientOptions function
	clientOptions := options.Client().ApplyURI(MongoDB)

	client, err := mongo.Connect(clientOptions)

	if err != nil{
		return nil
	}

	return client


 }

 func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load("/Users/kevindev/Desktop/movienightv2/backend/.env")

	if err != nil {
		log.Panicln("Warning: unable to find .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
 }