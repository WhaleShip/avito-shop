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
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(16) PRIMARY KEY, 
			password CHAR(64) NOT NULL,
			coins BIGINT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS coin_transactions (
			id SERIAL PRIMARY KEY,
			from_user VARCHAR(16) REFERENCES users(username) ON DELETE CASCADE,
			to_user VARCHAR(16) REFERENCES users(username) ON DELETE CASCADE,
			amount BIGINT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_coin_transactions_from ON coin_transactions(from_user);`,
		`CREATE INDEX IF NOT EXISTS idx_coin_transactions_to ON coin_transactions(to_user);`,
		`CREATE TABLE IF NOT EXISTS inventory_items (
			id SERIAL PRIMARY KEY,
			user_username VARCHAR(16) REFERENCES users(username) ON DELETE CASCADE,
			item_name VARCHAR(16) NOT NULL,
			quantity INT NOT NULL DEFAULT 0,
			UNIQUE (user_username, item_name)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_inventory_items_user ON inventory_items(user_username);`,
		`CREATE TABLE IF NOT EXISTS merch_items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(16) NOT NULL UNIQUE,
			price BIGINT NOT NULL
		);`,
		`INSERT INTO merch_items (name, price) VALUES 
			('t-shirt', 80),
			('cup', 20),
			('book', 50),
			('pen', 10),
			('powerbank', 200),
			('hoody', 300),
			('umbrella', 200),
			('socks', 10),
			('wallet', 50),
			('pink-hoody', 500)
		ON CONFLICT DO NOTHING;`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(context.Background(), q); err != nil {
			return fmt.Errorf("migration error: %w", err)
		}
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
		SigningKey:   jwtware.SigningKey{Key: *config.GetJWTSecret()},
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
