//go:build integration
// +build integration

package integration_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/whaleship/avito-shop/internal/config"
	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/handlers"
	"github.com/whaleship/avito-shop/internal/utils"
)

var testApp *fiber.App
var dbContainer tc.Container

func runMigrations(pool *pgxpool.Pool) error {
	data, err := os.ReadFile("../../migrations/init.sql")
	if err != nil {
		return fmt.Errorf("reading file error: %w", err)
	}

	if _, err := pool.Exec(context.Background(), string(data)); err != nil {
		fmt.Errorf("migration error: %w", err)
	}

	return nil
}

func NewApp() *fiber.App {
	if err := godotenv.Load(); err != nil {
		log.Print("error loading .env file: ", err)
	}

	app := fiber.New()

	pool, err := database.ConnectPostgresPool()
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", pool)
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

	return app
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := tc.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		fmt.Println("Failed to start container:", err)
		os.Exit(1)
	}
	dbContainer = container

	host, err := container.Host(ctx)
	if err != nil {
		fmt.Println("Failed to get container host:", err)
		os.Exit(1)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		fmt.Println("Failed to get container port:", err)
		os.Exit(1)
	}

	os.Setenv("PGBOUNCER_HOST", host)
	os.Setenv("PGBOUNCER_PORT", port.Port())
	os.Setenv("POSTGRES_USER", "testuser")
	os.Setenv("POSTGRES_PASSWORD", "testpass")
	os.Setenv("POSTGRES_DB", "testdb")
	os.Setenv("SSL_MODE", "disable")

	pool, err := database.ConnectPostgresPool()
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	if err := runMigrations(pool); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	testApp = NewApp()

	code := m.Run()

	if err := container.Terminate(ctx); err != nil {
		fmt.Println("Failed to terminate container:", err)
	}
	os.Exit(code)
}
