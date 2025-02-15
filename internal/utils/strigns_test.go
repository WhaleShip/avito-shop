package utils

import "testing"

// TestCapitalizeFirst проверяет функцию CapitalizeFirst.
func TestCapitalizeFirst(t *testing.T) {
	// Пустая строка
	if got := CapitalizeFirst(""); got != "" {
		t.Errorf("ожидалась пустая строка, получено '%s'", got)
	}
	// Строка с маленькой буквы
	if got := CapitalizeFirst("hello"); got != "Hello" {
		t.Errorf("ожидалось 'Hello', получено '%s'", got)
	}
	// Строка с уже заглавной буквы
	if got := CapitalizeFirst("Hello"); got != "Hello" {
		t.Errorf("ожидалось 'Hello', получено '%s'", got)
	}
	// Строка на кириллице
	if got := CapitalizeFirst("привет"); got != "Привет" {
		t.Errorf("ожидалось 'Привет', получено '%s'", got)
	}
}
