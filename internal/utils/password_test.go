package utils

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mypassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Error("неожиданная ошибка при хэшировании: ", err)
	}
	if hashed == password || len(hashed) == 0 {
		t.Error("хэш пароля не должен совпадать с исходным паролем и быть пустым")
	}

	if !CheckPassword(hashed, password) {
		t.Error("ожидалось успешное сравнение пароля")
	}

	if CheckPassword(hashed, "wrongpassword") {
		t.Error("ожидалось, что неправильный пароль не пройдет проверку")
	}
}
