package utils

import "testing"

func TestCapitalizeFirst(t *testing.T) {
	if got := CapitalizeFirst(""); got != "" {
		t.Error("ожидалась пустая строка, получено ", got)
	}
	if got := CapitalizeFirst("hello"); got != "Hello" {
		t.Error("ожидалось 'Hello', получено ", got)
	}
	if got := CapitalizeFirst("Hello"); got != "Hello" {
		t.Error("ожидалось 'Hello', получено ", got)
	}
	if got := CapitalizeFirst("привет"); got != "Привет" {
		t.Error("ожидалось 'Привет', получено ", got)
	}
}
