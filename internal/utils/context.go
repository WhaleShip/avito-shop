package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/whaleship/avito-shop/internal/database"
)

func ExtractDB(c *fiber.Ctx) (database.PgxIface, error) {
	dbConn := c.Locals("db")
	if dbConn == nil {
		return nil, errors.New("db connection not found in context")
	}
	db, ok := dbConn.(database.PgxIface)
	if !ok {
		return nil, errors.New("db connection type assertion failed")
	}
	return db, nil
}
