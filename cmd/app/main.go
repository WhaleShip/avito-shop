package main

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: ", err)
	}

	session, err := database.GetInitializedDb()
	if err != nil {
		log.Fatal("error connection DB: ", err)
	}
	defer log.Println(session.Close(context.Background()))

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", session)
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":8080"))
}
