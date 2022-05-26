package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/rycastr/simple-api/handler"
	repo "github.com/rycastr/simple-api/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	repo.Mongo = &repo.MongoInstance{Client: client}

	if dbName := os.Getenv("MONGODB_DB"); dbName != "" {
		repo.Mongo.Database = repo.Mongo.Client.Database(dbName)
	}

	// Fiber app
	app := fiber.New()

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/", handler.SignUp)
	auth.Post("/signin", handler.SignIn)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatal(err)
		}
	}()

	// gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Shutting down...")
	app.Shutdown()
	repo.Mongo.Client.Disconnect(context.TODO())
}
