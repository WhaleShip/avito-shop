// _test.go
package utils

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/valyala/fasthttp"
)

func TestExtractDB(t *testing.T) {
	t.Run("Нет соединения с БД в контексте", func(t *testing.T) {
		app := fiber.New()
		// Создаём пустой контекст fasthttp
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		db, err := ExtractDB(c)
		if err == nil {
			t.Error("ожидалась ошибка, получено nil")
		}
		if db != nil {
			t.Error("ожидался nil для db")
		}
	})

	t.Run("Неверный тип соединения", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		// Передаём значение неверного типа
		c.Locals("db", "не БД")
		db, err := ExtractDB(c)
		if err == nil || err.Error() != "db connection type assertion failed" {
			t.Errorf("ожидалась ошибка 'db connection type assertion failed', получено: %v", err)
		}
		if db != nil {
			t.Error("ожидался nil для db")
		}
	})

	t.Run("Успешное получение соединения", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		// Создаём фиктивное соединение pgx
		dummyConn := new(pgx.Conn)
		c.Locals("db", dummyConn)
		db, err := ExtractDB(c)
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if db != dummyConn {
			t.Error("ожидалось получение того же объекта соединения")
		}
	})
}

// TestGetUsername проверяет функцию GetUsername.
