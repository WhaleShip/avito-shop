//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/whaleship/avito-shop/internal/dto"
)

func TestBuyMerchEndpoint(t *testing.T) {
	t.Run("Успешная покупка предмета", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "buyer",
			"password": "buy123",
		}
		body, _ := json.Marshal(reqBody)
		reqAuth := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		reqAuth.Header.Set("Content-Type", "application/json")
		respAuth, err := testApp.Test(reqAuth, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		var authResp map[string]string
		err = json.NewDecoder(respAuth.Body).Decode(&authResp)
		assert.NoError(t, err)
		token := authResp["token"]

		reqBuy := httptest.NewRequest("GET", "/api/buy/t-shirt", nil)
		reqBuy.Header.Set("Authorization", "Bearer "+token)
		respBuy, err := testApp.Test(reqBuy, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, respBuy.StatusCode)

		reqInfo := httptest.NewRequest("GET", "/api/info", nil)
		reqInfo.Header.Set("Authorization", "Bearer "+token)
		respInfo, err := testApp.Test(reqInfo, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, respInfo.StatusCode)
		var info dto.InfoResponse
		bodyBytes, _ := ioutil.ReadAll(respInfo.Body)
		err = json.Unmarshal(bodyBytes, &info)
		assert.NoError(t, err)
		found := false
		for _, item := range info.Inventory {
			if item.Type == "t-shirt" && item.Quantity > 0 {
				found = true
				break
			}
		}
		assert.True(t, found, "t-shirt должен быть в инвентаре")
	})

	t.Run("Покупка несуществующего предмета", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "buyer2",
			"password": "buy123",
		}
		body, _ := json.Marshal(reqBody)
		reqAuth := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(body))
		reqAuth.Header.Set("Content-Type", "application/json")
		respAuth, err := testApp.Test(reqAuth, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		var authResp map[string]string
		err = json.NewDecoder(respAuth.Body).Decode(&authResp)
		assert.NoError(t, err)
		token := authResp["token"]

		reqBuy := httptest.NewRequest("GET", "/api/buy/nonexistent-item", nil)
		reqBuy.Header.Set("Authorization", "Bearer "+token)
		respBuy, err := testApp.Test(reqBuy, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, respBuy.StatusCode)
	})

	t.Run("Покупка предмета без токена", func(t *testing.T) {
		reqBuy := httptest.NewRequest("GET", "/api/buy/t-shirt", nil)
		respBuy, err := testApp.Test(reqBuy, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, respBuy.StatusCode)
	})
}
