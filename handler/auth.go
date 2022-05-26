package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rycastr/simple-api/model"
	repo "github.com/rycastr/simple-api/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SignUp(c *fiber.Ctx) error {
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// force MongoDB to generate an ID
	user.ID = ""
	user.SetCredentials()

	result, err := repo.Mongo.Database.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	// convert the ObjectID to a string and set it as the response body
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return c.JSON(user)
}

func SignIn(c *fiber.Ctx) error {
	var userCredentials map[string]string
	if err := c.BodyParser(&userCredentials); err != nil {
		return err
	}

	filter := bson.M{"email": userCredentials["email"]}

	var user model.User

	result := repo.Mongo.Database.Collection("users").FindOne(context.TODO(), filter)
	if err := result.Decode(&user); err != nil {
		return err
	}

	if !user.CheckCredentials(userCredentials["password"]) {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	return c.JSON(user)
}
