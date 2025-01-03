package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	WarnAdminAccess      = "ВХОД В АДМИНКУ"
	WarnTooLong          = "СЛИШКОМ ДОЛГАЯ ОБРАБОТКА"
	InfoUpdateFulfilled  = "Обновление обработано"
	InfoRequestFulfilled = "Обработан запрос"
)

var (
	FromAdminPanel           = "AdminPanel"
	FromUnsupported          = "Unsupported"
	FromNew                  = "New"
	FromListShorties         = "ListShorties"
	FromShortyInfo           = "ShortyInfo"
	FromDeleteShorty         = "DeleteShorty"
	FromDeleteShortyPrompt   = "DeleteShortyPrompt"
	FromDeleteSelectedShorty = "DeleteSelectedShorty"
	FromMenu                 = "Menu"
)

var (
	OriginBot    = "BOT"
	OriginServer = "SERVER"
)

// Структура, хранящая переменную для доступа к логгеру. Также реализует Loggier.
type Logger struct {
	logger *zerolog.Logger
}

// Интерфейс, содержащий методы для работы с логгером.
type Loggier interface {
	InfoBot(msg string, origin string, from string, userID int64, username string, start time.Time)
	WarnBot(msg string, origin string, from string, userID int64, username string, start time.Time)
	ErrorBot(msg string, origin string)
	InfoServer(msg string, shortyURL string, originalURL string, start time.Time)
	WarnServer(msg string, shortyURL string, originalURL string, start time.Time)
	ErrorServer(msg string, code int)
}

// Инициализирует логгер, возвращает структуру *Logger.
func InitLogger() *Logger {
	logger := zerolog.New(&zerolog.ConsoleWriter{
		Out:          os.Stdout,
		TimeFormat:   "02.01.2006 | 15:04:05",
		TimeLocation: time.UTC,
	}).With().Timestamp().Logger()

	return &Logger{logger: &logger}
}

// Логгирует Info в боте.
func (l *Logger) InfoBot(msg string, origin string, from string, userID int64, username string, start time.Time) {
	if username == "" {
		username = "отсутствует"
	}

	if time.Since(start) >= time.Duration(time.Millisecond*10) {
		l.WarnBot(WarnTooLong, origin, from, userID, username, start)
	}

	l.logger.Info().
		Str("ORIGIN", origin).
		Str("FROM", from).
		Int64("USER_ID", userID).
		Str("USERNAME", username).
		TimeDiff("SPEED", time.Now(), start).
		Msg(msg)
}

// Логгирует Warning в боте.
func (l *Logger) WarnBot(msg string, origin string, from string, userID int64, username string, start time.Time) {
	if username == "" {
		username = "отсутствует"
	}

	l.logger.Warn().
		Str("ORIGIN", origin).
		Str("FROM", from).
		Int64("USER_ID", userID).
		Str("USERNAME", username).
		TimeDiff("SPEED", time.Now(), start).
		Msg(msg)
}

// Логгирует Error в боте.
func (l *Logger) ErrorBot(msg string, origin string) {
	l.logger.Error().
		Str("ORIGIN", origin).
		Msg(msg)
}

// Логгирует Info на сервере.
func (l *Logger) InfoServer(msg string, shortyURL string, originalURL string, start time.Time) {
	if time.Since(start) >= time.Duration(time.Second*2) {
		l.WarnServer(WarnTooLong, shortyURL, originalURL, start)
	}

	l.logger.Info().
		Str("SHORTY", shortyURL).
		Str("REDIRECT", originalURL).
		TimeDiff("SPEED", time.Now(), start).
		Msg(msg)
}

// Логгирует Warning на сервере.
func (l *Logger) WarnServer(msg string, shortyURL string, originalURL string, start time.Time) {
	l.logger.Warn().
		Str("SHORTY", shortyURL).
		Str("REDIRECT", originalURL).
		TimeDiff("SPEED", time.Now(), start).
		Msg(msg)
}

// Логгирует Error на сервере.
func (l *Logger) ErrorServer(msg string, code int) {
	l.logger.Error().
		Int("CODE", code).
		Msg(msg)
}
