# 🔗 shortybot

**shortybot** - это сервис сокращения ссылок в формате телеграмм-бота.

Простой и удобный интерфейс, очеловеченные и содержательные ответы, функционал, позволяющий осуществлять все необходимые действия - **shortybot**!

<img src="https://github.com/user-attachments/assets/37b7bb53-8e76-4d63-8728-3cccf64bf67f">

## Основные функции

- **Сокращение ссылок:**   
Просто отправьте ссылку боту, и он создаст сокращенную версию.

- **Хранение сокращенных ссылок:**   
Все ваши сокращенные ссылки автоматически сохраняются. Вы можете просмотреть их в разделе «Мои Shorties», где будет доступна информация о дате создания.

- **Удаление сокращенных ссылок:**   
Если сокращенная ссылка больше не нужна, ее можно удалить через раздел «Удалить Shorty».

- **Безопасность:**   
Сокращенные ссылки доступны только вам. Бот использует уникальный идентификатор, связанный с вашим Telegram-аккаунтом, для защиты данных.

---

## Детали реализации

**shortybot** состоит из двух основных компонентов:

**Бот:**   
Отвечает за взаимодействие с пользователем через Telegram, включая обработку запросов, создание ссылок и управление ими.

**Сервер:**   
Обрабатывает входящие GET-запросы с уникальными идентификаторами сокращенных ссылок, извлекает информацию из базы данных и перенаправляет пользователя на оригинальный URL.

### Как это работает?

1. Пользователь отправляет ссылку боту для сокращения.
2. Бот сохраняет запись в базе данных, включая:
- Уникальный идентификатор Telegram-пользователя.
- Оригинальную ссылку.
- Уникальный идентификатор сокращенной ссылки.
- Дату создания.
4. Пользователь переходит по сокращенной ссылке.
5. Сервер получает запрос, извлекает идентификатор, обращается к базе данных и выполняет редирект на оригинальный URL.

### Ограничения

- Нельзя создать две сокращенные ссылки для одной и той же оригинальной ссылки.
- Нельзя сократить уже сокращенную ссылку.
- Нельзя сократить нерабочую ссылку.

## Технологии

- [**Go**](https://go.dev/) - основной язык разработки.
- [**Fiber**](https://docs.gofiber.io/) - веб-фреймворк для сервера.
- [**Telebot**](https://github.com/tucnak/telebot) - библиотека для работы с Telegram-ботом.
- [**PostgreSQL**](https://www.postgresql.org/) - база данных.
- [**GORM**](https://gorm.io/) - ORM для взаимодействия с базой данных.
- [**Zerolog**](https://github.com/rs/zerolog) - для структурированного логгирования.
- [**Docker**](https://www.docker.com/) - для контейнеризации.

## Структура проекта

```
shortybot/
├── cmd/
│   ├── bot/
│   │   └── main.go
│   ├── server/
│   │   └── main.go
│   ├── bot-Dockerfile
│   ├── server-Dockerfile
├── internal/
│   ├── bot/
│   │   ├── app/
│   │   │   └── app.go
│   │   ├── handlers/
│   │   │   └── handlers.go
│   │   ├── helpers/
│   │   │   └── helpers.go
│   │   ├── middleware/
│   │   │   └── middleware.go
│   │   ├── models/
│   │       ├── buttons.go
│   │       └── responses.go
│   ├── server/
│   │   ├── app/
│   │   │   └── app.go
│   │   ├── handlers/
│   │   │   └── handlers.go
│   │   ├── models/
│   │       └── responses.go
│   ├── db/
│   │   └── db.go
│   ├── logger/
│       └── logger.go
├── .dockerignore
├── .gitignore
├── docker-compose.yaml
├── go.mod
└── go.sum
```

## Лицензия

Проект распространяется под лицензией [**MIT**](https://mit-license.org/). Вы можете свободно использовать и модифицировать код при соблюдении условий лицензии.