package models

// Модель ответа.
type Response struct {
	Message string  `json:"message"`
	Status  *Status `json:"status"`
}

// Модель ошибки.
type Error struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Status  *Status `json:"status"`
}

// Модель статуса.
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
