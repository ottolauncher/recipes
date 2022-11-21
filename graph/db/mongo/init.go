package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TODO add env support

const uri = "mongodb://127.0.0.1:27017/recipedb"

func Init() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(":::: DATABASE CONNECTION ERROR :::::", err)
		panic(err)
	}
	log.Println("Connected Successfully")
	return client
}
