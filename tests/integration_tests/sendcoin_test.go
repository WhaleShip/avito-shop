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

func TestSendCoinEndpoint(t *testing.T) {
	t.Run("Успешная отправка монет", func(t *testing.T) {
		senderBody, _ := json.Marshal(map[string]string{
			"username": "sender",
			"password": "pass123",
		})
		receiverBody, _ := json.Marshal(map[string]string{
			"username": "receiver",
			"password": "pass456",
		})

		reqSender := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(senderBody))
		reqSender.Header.Set("Content-Type", "application/json")
		respSender, err := testApp.Test(reqSender, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		var senderAuth map[string]string
		err = json.NewDecoder(respSender.Body).Decode(&senderAuth)
		assert.NoError(t, err)
		senderToken := senderAuth["token"]

		reqReceiver := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(receiverBody))
		reqReceiver.Header.Set("Content-Type", "application/json")
		_, err = testApp.Test(reqReceiver, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)

		sendBody, _ := json.Marshal(map[string]interface{}{
			"toUser": "receiver",
			"amount": 100,
		})
		reqSend := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewReader(sendBody))
		reqSend.Header.Set("Content-Type", "application/json")
		reqSend.Header.Set("Authorization", "Bearer "+senderToken)
		respSend, err := testApp.Test(reqSend, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, respSend.StatusCode)
	})

	t.Run("Отправка монет несуществующему пользователю", func(t *testing.T) {
		senderBody, _ := json.Marshal(map[string]string{
			"username": "sender2",
			"password": "pass123",
		})
		reqSender := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(senderBody))
		reqSender.Header.Set("Content-Type", "application/json")
		respSender, err := testApp.Test(reqSender, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		var senderAuth map[string]string
		err = json.NewDecoder(respSender.Body).Decode(&senderAuth)
		assert.NoError(t, err)
		senderToken := senderAuth["token"]

		sendBody, _ := json.Marshal(map[string]interface{}{
			"toUser": "nonexistent",
			"amount": 50,
		})
		reqSend := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewReader(sendBody))
		reqSend.Header.Set("Content-Type", "application/json")
		reqSend.Header.Set("Authorization", "Bearer "+senderToken)
		respSend, err := testApp.Test(reqSend, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, respSend.StatusCode)
	})

	t.Run("Отправка монет с некорректными данными", func(t *testing.T) {
		senderBody, _ := json.Marshal(map[string]string{
			"username": "sender3",
			"password": "pass123",
		})
		reqSender := httptest.NewRequest("POST", "/api/auth", bytes.NewReader(senderBody))
		reqSender.Header.Set("Content-Type", "application/json")
		respSender, err := testApp.Test(reqSender, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		var senderAuth map[string]string
		err = json.NewDecoder(respSender.Body).Decode(&senderAuth)
		assert.NoError(t, err)
		senderToken := senderAuth["token"]

		sendBody, _ := json.Marshal(map[string]interface{}{
			"amount": 50,
		})
		reqSend := httptest.NewRequest("POST", "/api/sendCoin", bytes.NewReader(sendBody))
		reqSend.Header.Set("Content-Type", "application/json")
		reqSend.Header.Set("Authorization", "Bearer "+senderToken)
		respSend, err := testApp.Test(reqSend, int(5*time.Second/time.Millisecond))
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, respSend.StatusCode)
	})
}
