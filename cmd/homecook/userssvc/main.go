package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/quocdat32461997/HomeCOOK/internal/cloud"
	"github.com/quocdat32461997/HomeCOOK/internal/services/users"
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

	// Initialize the user service server and proxy
	userService := &users.Server{
		Mongo:    mongo,
		Endpoint: "0.0.0.0:9000",
	}

	go users.StartUserService(userService)
	go users.StartUserServiceProxy(userService)

	// Wait for Control+C
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	log.Println("User's service waiting for SIGINT...")
	<-ch

}
