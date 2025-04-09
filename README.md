# 🔗 shortybot

**shortybot** is a link shortening service in the format of a Telegram bot.

A simple and convenient interface, human-friendly and informative responses, and functionality that allows all necessary actions — **shortybot**!

<img src="https://github.com/user-attachments/assets/37b7bb53-8e76-4d63-8728-3cccf64bf67f">

## Main Features

- **Link Shortening:**  
Simply send a link to the bot, and it will create a shortened version.

- **Storage of Shortened Links:**  
All your shortened links are automatically saved. You can view them in the "My Shorties" section, where creation date information will be available.

- **Deleting Shortened Links:**  
If a shortened link is no longer needed, it can be deleted via the "Delete Shorty" section.

- **Security:**  
Shortened links are accessible only to you. The bot uses a unique identifier linked to your Telegram account to protect your data.

## Implementation Details

**shortybot** consists of two main components:

**Bot:**  
Responsible for interacting with the user through Telegram, including processing requests, creating links, and managing them.

**Server:**  
Handles incoming GET requests with unique shortened link identifiers, retrieves information from the database, and redirects the user to the original URL.

### How does it work?

1. The user sends a link to the bot to shorten.  
2. The bot saves a record in the database, including:  
   - A unique Telegram user ID.  
   - The original link.  
   - A unique shortened link ID.  
   - The creation date.  
3. The user follows the shortened link.  
4. The server receives the request, extracts the identifier, queries the database, and redirects to the original URL.

### Limitations

- You cannot create two shortened links for the same original link.  
- You cannot shorten an already shortened link.  
- You cannot shorten a non-working link.

## Project Structure

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

## License

[**MIT**](https://mit-license.org/)
