package handlers

import (
	"os"
	"strings"

	"gopkg.in/telebot.v4"

	"gorm.io/gorm"

	"github.com/xoticdsign/shortybot/internal/bot/helpers"
	"github.com/xoticdsign/shortybot/internal/bot/models"
	"github.com/xoticdsign/shortybot/internal/db"
)

// Структура, хранящая все необходимые хендлерам бота зависимости.
type Dependencies struct {
	DB      db.Querier
	Helpers helpers.Helpers
}

// Отлавливает ошибки, обрабатывает и направляет соответствующее переменную в контексте в хенделер Menu.
func (d *Dependencies) OnError(err error, c telebot.Context) {
	c.Set("error", models.MsgOnError+err.Error())

	d.Menu(c)
}

// Обрабатывает обновления, неподдерживаемые ботом.
func (d *Dependencies) Unsupported(c telebot.Context) error {
	switch {
	case c.Text() == "/start":
		return d.Menu(c)

	case strings.Contains(c.Text(), "https://") || strings.Contains(c.Text(), "http://"):
		c.Set("url", c.Text())

		return d.New(c)

	default:
		c.Set("msg", models.FailedGlobalUnsupportedCmd+"\n\n")

		return d.Menu(c)
	}
}

// Добавляет новую сокращенную ссылку в БД после всех необходимых проверок.
func (d *Dependencies) New(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user := c.Get("user").(*telebot.User)
	url := c.Get("url").(string)

	ok := strings.Contains(url, publicAdr)
	if ok {
		return c.Send(models.FailedNewCantShortShorty, models.ReplyReturnToMenuWithError)
	}

	ok = d.Helpers.CheckURL(url)
	if !ok {
		return c.Send(models.FailedNewIncorrectURL, models.ReplyReturnToMenuWithError)
	}

	shorty := d.Helpers.ShortyGenerator(7)

	err := d.DB.New(user.Username, url, shorty)
	if err != nil {
		switch {
		case err == gorm.ErrDuplicatedKey:
			return c.Send(models.FailedNewDuplicate, models.ReplyReturnToMenuWithError)

		case err == gorm.ErrCheckConstraintViolated:
			return c.Send(models.FailedNewLimitExceeded, models.ReplyReturnToMenuWithError)

		default:
			return err
		}
	}
	c.Set("msg", models.SuccessNew+publicAdr+shorty+"\n\n")

	return d.Menu(c)
}

// Отвечает за работу кнопки "Мои Shorties".
func (d *Dependencies) ListShorties(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user := c.Get("user").(*telebot.User)

	shorties, err := d.DB.ListShorties(user.Username)
	if err != nil {
		switch {
		case err == gorm.ErrRecordNotFound:
			return c.EditOrSend(models.FailedGlobalNoShorties, models.ReplyReturnToMenu)

		default:
			return err
		}
	}
	var btns [][]telebot.InlineButton

	for _, shorty := range shorties {
		var btn []telebot.InlineButton

		btnShortyInfo := telebot.InlineButton{
			Unique: models.BtnShortyInfo.Unique,
			Text:   publicAdr + shorty.ShortyURL,
			Data:   shorty.ShortyURL,
		}

		btn = append(btn, btnShortyInfo)
		btns = append(btns, btn)
	}

	var btnReturnToMenu []telebot.InlineButton

	btnReturnToMenu = append(btnReturnToMenu, models.BtnReturnToMenu)
	btns = append(btns, btnReturnToMenu)

	return c.EditOrSend(models.MsgListShorties, &telebot.ReplyMarkup{
		InlineKeyboard: btns,
	})
}

// Отвечает за работу кнопок, созданных в результате работы функции ListShorties.
func (d *Dependencies) ShortyInfo(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	shortyURL := c.Data()

	shorty, err := d.DB.ShortyInfo(shortyURL)
	if err != nil {
		switch {
		case err == gorm.ErrRecordNotFound:
			return c.EditOrSend(models.FailedGlobalNoShorties, models.ReplyReturnToListShorties)

		default:
			return err
		}
	}
	shortyInfo := []string{
		models.MsgShortyInfo,
		"· Shorty:\n" + publicAdr + shorty.ShortyURL,
		"· Оригинальная ссылка:\n" + shorty.OriginalURL,
		"· Дата создания:\n" + shorty.DateCreated.Format("02.01.2006, 15:04"),
	}

	shortyInfoFmt := strings.Join(shortyInfo, "\n\n")

	return c.EditOrSend(shortyInfoFmt, models.ReplyReturnToListShorties)
}

// Отвечает за работу кнопки "Удалить Shorty".
func (d *Dependencies) DeleteShorty(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user := c.Get("user").(*telebot.User)

	shorties, err := d.DB.ListShorties(user.Username)
	if err != nil {
		switch {
		case err == gorm.ErrRecordNotFound:
			return c.EditOrSend(models.FailedGlobalNoShorties, models.ReplyReturnToMenu)

		default:
			return err
		}
	}
	var btns [][]telebot.InlineButton

	for _, shorty := range shorties {
		var btn []telebot.InlineButton

		btnDeleteShortyPrompt := telebot.InlineButton{
			Unique: models.BtnDeleteShortyPrompt.Unique,
			Text:   publicAdr + shorty.ShortyURL,
			Data:   shorty.ShortyURL,
		}

		btn = append(btn, btnDeleteShortyPrompt)
		btns = append(btns, btn)
	}

	var btnReturnToMenu []telebot.InlineButton

	btnReturnToMenu = append(btnReturnToMenu, models.BtnReturnToMenu)
	btns = append(btns, btnReturnToMenu)

	return c.EditOrSend(models.MsgDeleteShorty, &telebot.ReplyMarkup{
		InlineKeyboard: btns,
	})
}

// Отвечает за работу кнопок "Да" и "Нет" при действии промпта на удаление ссылки.
func (d *Dependencies) DeleteShortyPrompt(c telebot.Context) error {
	shortyURL := c.Data()

	return c.EditOrSend(models.MsgDeleteShortyPrompt, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{
				{
					Unique: models.BtnDeleteSelectedShorty.Unique,
					Text:   models.BtnDeleteSelectedShorty.Text,
					Data:   shortyURL,
				},
				{
					Unique: models.BtnDeleteSelectedShorty2.Unique,
					Text:   models.BtnDeleteSelectedShorty2.Text,
				},
			},
		},
	})
}

// Удаляет сокращенную ссылку из БД.
func (d *Dependencies) DeleteSelectedShorty(c telebot.Context) error {
	shortyURL := c.Data()

	err := d.DB.DeleteShorty(shortyURL)
	if err != nil {
		switch {
		case err == gorm.ErrRecordNotFound:
			return c.EditOrSend(models.FailedGlobalNoShorties, models.ReplyReturnToListShorties)

		default:
			return err
		}
	}
	c.Set("msg", models.SuccessDelete+"\n\n")

	return d.Menu(c)
}

// Отвечает за работу главного меню.
func (d *Dependencies) Menu(c telebot.Context) error {
	var data string

	if c.Callback() != nil {
		data = c.Callback().Data
	}

	if data == "failed" {
		return c.Send(models.MsgMenuGreeting+models.MsgMenuDescrtiption, models.ReplyMenu)
	}

	msgErr, ok := c.Get("error").(string)
	if ok {
		return c.Send(msgErr+models.MsgMenuDescrtiption, models.ReplyMenu)
	}

	msg, ok := c.Get("msg").(string)
	if ok {
		return c.EditOrSend(msg+models.MsgMenuDescrtiption, models.ReplyMenu)
	}
	return c.EditOrSend(models.MsgMenuGreeting+models.MsgMenuDescrtiption, models.ReplyMenu)
}
