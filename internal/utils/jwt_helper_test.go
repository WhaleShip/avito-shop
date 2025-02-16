package utils

import (
	"encoding/json"
	"errors"
	"strconv"
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
			t.Error("неверное сообщение об ошибке: ", err.Error())
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
		token.Claims = jwt.RegisteredClaims{}
		c.Locals("user", token)
		username, err := GetUsername(c)
		if err == nil || username != "" {
			t.Errorf("ожидалась ошибка при неверном типе claims, получено: %s, %v", username, err)
		}
		if err.Error() != "failed to convert token to MapClaims" {
			t.Error("неверное сообщение об ошибке: ", err.Error())
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
			t.Error("неверное сообщение об ошибке: ", err.Error())
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
			t.Error("неожиданная ошибка: ", err)
		}
		if username != expectedUsername {
			t.Errorf("ожидалось username '%s', получено '%s'", expectedUsername, username)
		}
	})
}

func TestJwtError(t *testing.T) {
	app := fiber.New()
	reqCtx := new(fasthttp.RequestCtx)
	c := app.AcquireCtx(reqCtx)
	defer app.ReleaseCtx(c)

	err := JwtError(c, errors.New("какая-то ошибка"))
	if err != nil {
		t.Error("неожиданная ошибка: ", err)
	}

	if c.Response().StatusCode() != fiber.StatusUnauthorized {
		t.Error("ожидался статус " + strconv.Itoa(fiber.StatusUnauthorized) +
			", получен " + strconv.Itoa(c.Response().StatusCode()))
	}

	var body map[string]string
	if err := json.Unmarshal(c.Response().Body(), &body); err != nil {
		t.Error("ошибка при разборе JSON: ", err)
	}

	if body["errors"] != "Неавторизован" {
		t.Error("ожидалось сообщение Неавторизован, получено ", body["errors"])
	}
}
