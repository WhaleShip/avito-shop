package main

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/config"
	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/handlers"
	"github.com/whaleship/avito-shop/internal/utils"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("error loading .env file: ", err)
	}

	session, err := database.GetInitializedDb()
	if err != nil {
		log.Fatal("error connection DB: ", err)
	}
	defer func() {
		err := session.Close(context.Background())
		log.Println("Closing connection:", err)
	}()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", session)
		return c.Next()
	})

	app.Post("/api/auth", handlers.AuthHandler)

	app.Use("/api", jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.GetJWTSecret())},
		ErrorHandler: utils.JwtError,
	}))

	app.Get("/api/info", handlers.InfoHandler)
	app.Get("/api/buy/:item", handlers.BuyHandler)

	log.Fatal(app.Listen(":8080"))
}
