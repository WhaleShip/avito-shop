//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpoint(t *testing.T) {
	t.Run("Успешная аутентификация", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var authResp map[string]string
		err = json.NewDecoder(resp.Body).Decode(&authResp)
		assert.NoError(t, err)
		token, ok := authResp["token"]
		assert.True(t, ok, "Token должен присутствовать в ответе")
		assert.NotEmpty(t, token, "Token не должен быть пустым")
	})

	t.Run("Аутентификация с пустым паролем", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser_empty",
			"password": "",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestAuthInvalidPassword(t *testing.T) {
	t.Run("Аутентификация с неверным паролем", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser_invalid",
			"password": "correctpass",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		_, err := testApp.Test(req, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)

		reqBody["password"] = "wrongpass"
		body, _ = json.Marshal(reqBody)
		req = httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Аутентификация нового пользователя", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "nonexistent_user",
			"password": "somepass",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
