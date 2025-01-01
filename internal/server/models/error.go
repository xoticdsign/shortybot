package models

// Модель ошибки.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
