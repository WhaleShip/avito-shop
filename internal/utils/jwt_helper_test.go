package utils

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

func TestGetUsername(t *testing.T) {
	t.Run("Нет токена пользователя", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		username, err := GetUsername(c)
		if err == nil || username != "" {
			t.Errorf("ожидалась ошибка и пустой username, получено: %s, %v", username, err)
		}
		if err.Error() != "faild to get username" {
			t.Errorf("неверное сообщение об ошибке: %s", err.Error())
		}
	})

	t.Run("Неверный тип токена", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		c.Locals("user", "не токен")
		username, err := GetUsername(c)
		if err == nil || username != "" {
			t.Errorf("ожидалась ошибка при неверном типе токена, получено: %s, %v", username, err)
		}
	})

	t.Run("Неверный тип claims", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		token := jwt.New(jwt.SigningMethodHS256)
		// Устанавливаем claims неверного типа, который не является jwt.MapClaims
		token.Claims = jwt.RegisteredClaims{}
		c.Locals("user", token)
		username, err := GetUsername(c)
		if err == nil || username != "" {
			t.Errorf("ожидалась ошибка при неверном типе claims, получено: %s, %v", username, err)
		}
		if err.Error() != "failed to convert token to MapClaims" {
			t.Errorf("неверное сообщение об ошибке: %s", err.Error())
		}
	})

	t.Run("username не строка", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims = jwt.MapClaims{"username": 123} // не строка
		c.Locals("user", token)
		username, err := GetUsername(c)
		if err == nil || username != "" {
			t.Errorf("ожидалась ошибка при неверном типе username, получено: %s, %v", username, err)
		}
		if err.Error() != "failed to convert username to string" {
			t.Errorf("неверное сообщение об ошибке: %s", err.Error())
		}
	})

	t.Run("Успешное получение username", func(t *testing.T) {
		app := fiber.New()
		reqCtx := new(fasthttp.RequestCtx)
		c := app.AcquireCtx(reqCtx)
		defer app.ReleaseCtx(c)

		expectedUsername := "testuser"
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims = jwt.MapClaims{"username": expectedUsername}
		c.Locals("user", token)
		username, err := GetUsername(c)
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if username != expectedUsername {
			t.Errorf("ожидалось username '%s', получено '%s'", expectedUsername, username)
		}
	})
}

// TestJwtError проверяет функцию JwtError.
func TestJwtError(t *testing.T) {
	app := fiber.New()
	reqCtx := new(fasthttp.RequestCtx)
	c := app.AcquireCtx(reqCtx)
	defer app.ReleaseCtx(c)

	err := JwtError(c, errors.New("какая-то ошибка"))
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if c.Response().StatusCode() != fiber.StatusUnauthorized {
		t.Errorf("ожидался статус %d, получен %d", fiber.StatusUnauthorized, c.Response().StatusCode())
	}

	var body map[string]string
	if err := json.Unmarshal(c.Response().Body(), &body); err != nil {
		t.Errorf("ошибка при разборе JSON: %v", err)
	}

	if body["errors"] != "Неавторизован" {
		t.Errorf("ожидалось сообщение 'Неавторизован', получено '%s'", body["errors"])
	}
}
