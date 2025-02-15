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
	"github.com/whaleship/avito-shop/internal/dto"
)

func TestInfoEndpoint(t *testing.T) {
	t.Run("Успешное получение информации", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "testuser_info",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		reqAuth := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		reqAuth.Header.Set("Content-Type", "application/json")
		respAuth, err := testApp.Test(reqAuth, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, respAuth.StatusCode)

		var authResp map[string]string
		err = json.NewDecoder(respAuth.Body).Decode(&authResp)
		assert.NoError(t, err)
		token := authResp["token"]

		reqInfo := httptest.NewRequest("GET", "/api/info", nil)
		reqInfo.Header.Set("Authorization", "Bearer "+token)
		respInfo, err := testApp.Test(reqInfo, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, respInfo.StatusCode)

		var info dto.InfoResponse
		err = json.NewDecoder(respInfo.Body).Decode(&info)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, info.Coins, int64(0), "Количество монет не может быть отрицательным")
	})

	t.Run("Запрос без токена", func(t *testing.T) {
		reqInfo := httptest.NewRequest("GET", "/api/info", nil)
		respInfo, err := testApp.Test(reqInfo, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, respInfo.StatusCode)
	})

	t.Run("Запрос с некорректным токеном", func(t *testing.T) {
		reqInfo := httptest.NewRequest("GET", "/api/info", nil)
		reqInfo.Header.Set("Authorization", "Bearer invalidtoken")
		respInfo, err := testApp.Test(reqInfo, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, respInfo.StatusCode)
	})
}
