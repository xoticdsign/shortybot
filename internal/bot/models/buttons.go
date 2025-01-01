package models

import "gopkg.in/telebot.v4"

var (
	// Кнопка "Мои Shorties". Появляется в главном меню.
	BtnListShorties = telebot.InlineButton{
		Unique: "listShorties",
		Text:   "Мои Shorties",
	}

	// Кнопка "К списку Shorties". Появляется при отображение списка сокращенных ссылок.
	BtnReturnToListShorties = telebot.InlineButton{
		Unique: "listShorties",
		Text:   "К списку Shorties",
	}

	// Кнопка "здесь сокращенная ссылка". Появляется при отображение списка сокращенных ссылок.
	BtnShortyInfo = telebot.InlineButton{
		Unique: "shortyInfo",
	}
)

var (
	// Кнопка "Удалить Shorty". Появляется в главном меню.
	BtnDeleteShorty = telebot.InlineButton{
		Unique: "deleteShorty",
		Text:   "Удалить Shorty",
	}

	// Кнопка "здесь сокращенная ссылка". Появляется при отображение списка сокращенных ссылок в меню удаления.
	BtnDeleteShortyPrompt = telebot.InlineButton{
		Unique: "deleteShortyPrompt",
	}

	// Кнопка "Да". Появляется после при предупреждении перед удалением сокращенной ссылки.
	BtnDeleteSelectedShorty = telebot.InlineButton{
		Unique: "deleteSelectedShorty",
		Text:   "Да",
	}

	// Кнопка "Нет". Появляется после при предупреждении перед удалением сокращенной ссылки.
	BtnDeleteSelectedShorty2 = telebot.InlineButton{
		Unique: "deleteShorty",
		Text:   "Нет",
	}
)

var (
	// Кнопка "Вернуться в меню". Появляется, как предложение вернуться в меню на разных этапах взаимодействия.
	BtnReturnToMenu = telebot.InlineButton{
		Unique: "menu",
		Text:   "Вернуться в меню",
	}
)

var (
	// Набор кнопок, отображаемый, как предложение вернуться к списку сокращенных ссылок.
	ReplyReturnToListShorties = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{BtnReturnToListShorties},
		},
	}
)

var (
	// Набор кнопок, отображаемый в главном меню.
	ReplyMenu = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{BtnListShorties},
			{BtnDeleteShorty},
		},
	}

	// Набор кнопок, отображаемый, как предложение вернуться в главное меню.
	ReplyReturnToMenu = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{BtnReturnToMenu},
		},
	}

	// Набор кнопок, отображаемый, как предложение вернуться в главное меню. Дополнительно имеет при себе данные, которые свидетельствуют о том, что пользователь решил вернуться в главное меню после какого-то неудачного действия.
	ReplyReturnToMenuWithError = &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*BtnReturnToMenu.With("failed")},
		},
	}
)
