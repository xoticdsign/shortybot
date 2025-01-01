package helpers

import (
	"math/rand"
	"net/http"
	"time"
)

// Структура, к которой привязаны все вспомогательные функции.
type Helpers struct {
}

// Проверяет ссылку на действительность.
func (h *Helpers) CheckURL(url string) bool {
	client := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second * 10,
	}

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return false
	}
	return true
}

// Генерирует уникальный семизначный идентификатор сокращенной ссылки.
func (h *Helpers) ShortyGenerator(length int) string {
	const allowed = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)

	for i := 0; i < length; i++ {
		b[i] = allowed[random.Intn(len(allowed))]
	}

	return string(b)
}
