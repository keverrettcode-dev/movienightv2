package database

import (
	"log"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func DBInstance() *mongo.Client{ 
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: unable to find .env file!")
	}
 }