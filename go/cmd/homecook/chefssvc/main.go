package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/quocdat32461997/HomeCOOK/internal/cloud"
	"github.com/quocdat32461997/HomeCOOK/internal/services/chefs"
)

func main() {
	/*
		To Do:
		1. Replace `panic`s with log writes
	*/

	// Load environment
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Connect to MongoDB Atlas
	mongo := &cloud.MongoConn{
		Host:       os.Getenv("MONGO_HOST"),
		Authorizer: os.Getenv("MONGO_AUTH_DATABASE"),
		Database:   os.Getenv("MONGO_DATABASE"),
		Username:   os.Getenv("MONGO_USERNAME"),
		Password:   os.Getenv("MONGO_PASSWORD"),
	}
	session, err := mongo.InitDB()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Initialize the chef service server and proxy
	chefService := &chefs.Server{
		Mongo:    mongo,
		Endpoint: "0.0.0.0:9001",
	}
	go chefs.StartChefService(chefService)
	go chefs.StartChefServiceProxy(chefService)

	// Wait for Control+C
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	log.Println("Chef's service waiting for SIGINT...")
	<-ch

}
