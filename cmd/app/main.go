package main

import (
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
	if err := godotenv.Load(); err != nil {
		log.Print("error loading .env file: ", err)
	}

	dbPool, err := database.ConnectPostgresPool()
	if err != nil {
		log.Fatal("error connecting to DB: ", err)
	}
	defer func() {
		dbPool.Close()
		log.Println("Closing connection pool")
	}()

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", dbPool)
		return c.Next()
	})

	app.Post("/api/auth", handlers.AuthHandler)

	app.Use("/api", jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: config.GetJWTSecret()},
		ErrorHandler: utils.JwtError,
	}))

	app.Get("/api/info", handlers.InfoHandler)
	app.Get("/api/buy/:item", handlers.BuyHandler)
	app.Post("/api/sendCoin", handlers.SendCoinHandler)

	log.Fatal(app.Listen(":8080"))
}
