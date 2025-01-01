package handlers

import (
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v4"

	"gorm.io/gorm"

	"github.com/xoticdsign/shortybot/internal/bot/helpers"
	"github.com/xoticdsign/shortybot/internal/bot/models"
	"github.com/xoticdsign/shortybot/internal/db"
	"github.com/xoticdsign/shortybot/internal/logger"
)

// Структура, хранящая все необходимые хендлерам бота зависимости.
type Dependencies struct {
	DB      db.Querier
	Logger  logger.Loggier
	Helpers helpers.Helpers
}

// Отлавливает ошибки и логгирует.
func (d *Dependencies) OnError(err error, c telebot.Context) {
	d.Logger.ErrorBot(
		err.Error(),
		logger.OriginBot,
	)
}

// Админский хендлер. Отправляет админскую панель.
func (d *Dependencies) AdminPanel(c telebot.Context) error {
	var data string

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	_, ok = c.Get("admin").(string)
	if !ok {
		c.Delete()

		return c.Send(models.FailedGlobalUnsupportedCmd, models.ReplyReturnToMenuWithSend)
	}

	if c.Callback() != nil {
		data = c.Callback().Data
	}

	d.Logger.WarnBot(
		logger.WarnAdminAccess,
		logger.OriginBot,
		logger.FromAdminPanel,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	switch {
	case data == "send":
		if user.FirstName != "" {
			return c.Send(models.MsgAdminPanel1+user.FirstName+models.MsgAdminPanel2, models.ReplyAdminPanel)
		}
		return c.Send(models.MsgAdminPanel, models.ReplyAdminPanel)

	default:
		if user.FirstName != "" {
			return c.EditOrSend(models.MsgAdminPanel1+user.FirstName+models.MsgAdminPanel2, models.ReplyAdminPanel)
		}
		return c.EditOrSend(models.MsgAdminPanel, models.ReplyAdminPanel)
	}
}

// Админский хендлер. Отправляет количество уникальных пользователей и сокращенных ссылок.
func (d *Dependencies) AdminUsersAndShorties(c telebot.Context) error {
	_, ok := c.Get("admin").(string)
	if !ok {
		c.Delete()

		return c.Send(models.FailedGlobalUnsupportedCmd, models.ReplyReturnToMenuWithSend)
	}

	usersCount, shortiesCount, err := d.DB.UsersAndShorties()
	if err != nil {
		return c.EditOrSend(models.MsgAdminSuccess+"\n\n"+err.Error(), models.ReplyReturnToAdminPanel)
	}
	result := []string{
		models.MsgAdminSuccess,
		"· Пользователей:\n" + strconv.Itoa(usersCount),
		"· Ссылок:\n" + strconv.Itoa(int(shortiesCount)),
	}

	resultFmt := strings.Join(result, "\n\n")

	return c.EditOrSend(resultFmt, models.ReplyReturnToAdminPanel)
}

// Обрабатывает обновления, неподдерживаемые ботом.
func (d *Dependencies) Unsupported(c telebot.Context) error {
	isAdmin := false

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	_, ok = c.Get("admin").(string)
	if ok {
		isAdmin = true
	}

	switch {
	case c.Text() == "/start":
		return d.Menu(c)

	case c.Text() == "/admin" && isAdmin:
		return d.AdminPanel(c)

	case strings.Contains(c.Text(), "https://") || strings.Contains(c.Text(), "http://"):
		c.Set("url", c.Text())

		return d.New(c)

	default:
		d.Logger.InfoBot(
			logger.InfoUpdateFulfilled,
			logger.OriginBot,
			logger.FromUnsupported,
			user.ID,
			user.Username,
			c.Get("start").(time.Time),
		)

		return c.Send(models.FailedGlobalUnsupportedCmd, models.ReplyReturnToMenuWithSend)
	}
}

// Добавляет новую сокращенную ссылку в БД после всех необходимых проверок.
func (d *Dependencies) New(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	url, ok := c.Get("url").(string)
	if !ok {
		return telebot.ErrBadContext
	}

	ok = strings.Contains(url, publicAdr)
	if ok {
		return c.Send(models.FailedNewCantShortShorty, models.ReplyReturnToMenuWithSend)
	}

	ok = d.Helpers.CheckURL(url)
	if !ok {
		return c.Send(models.FailedNewIncorrectURL, models.ReplyReturnToMenuWithSend)
	}

	shorty := d.Helpers.ShortyGenerator(7)

	err := d.DB.New(user.ID, url, shorty)
	if err != nil {
		switch {
		case err == gorm.ErrDuplicatedKey:
			return c.Send(models.FailedNewDuplicate, models.ReplyReturnToMenuWithSend)

		case err == gorm.ErrCheckConstraintViolated:
			return c.Send(models.FailedNewLimitExceeded, models.ReplyReturnToMenuWithSend)

		default:
			return err
		}
	}

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromNew,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	return c.Send(models.SuccessNew+"\n\n"+publicAdr+shorty, models.ReplyReturnToMenuWithSend)
}

// Отвечает за работу кнопки "Мои Shorties".
func (d *Dependencies) ListShorties(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	shorties, err := d.DB.ListShorties(user.ID)
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

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromListShorties,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	return c.EditOrSend(models.MsgListShorties, &telebot.ReplyMarkup{
		InlineKeyboard: btns,
	})
}

// Отвечает за работу кнопок, созданных в результате работы функции ListShorties.
func (d *Dependencies) ShortyInfo(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

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

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromShortyInfo,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	return c.EditOrSend(shortyInfoFmt, models.ReplyReturnToListShorties)
}

// Отвечает за работу кнопки "Удалить Shorty".
func (d *Dependencies) DeleteShorty(c telebot.Context) error {
	publicAdr := os.Getenv("SERVER_PUBLIC_ADR")

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	shorties, err := d.DB.ListShorties(user.ID)
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

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromDeleteShorty,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	return c.EditOrSend(models.MsgDeleteShorty, &telebot.ReplyMarkup{
		InlineKeyboard: btns,
	})
}

// Отвечает за работу кнопок "Да" и "Нет" при действии промпта на удаление ссылки.
func (d *Dependencies) DeleteShortyPrompt(c telebot.Context) error {
	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	shortyURL := c.Data()

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromDeleteShortyPrompt,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

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
	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

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
	c.Set("msg", models.SuccessDelete)

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromDeleteSelectedShorty,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	return d.Menu(c)
}

// Отвечает за работу главного меню.
func (d *Dependencies) Menu(c telebot.Context) error {
	var data string
	var msgText string

	user, ok := c.Get("user").(*telebot.User)
	if !ok {
		return telebot.ErrBadContext
	}

	msg, ok := c.Get("msg").(string)
	if ok {
		msgText = msg
	}

	if c.Callback() != nil {
		data = c.Callback().Data
	}

	d.Logger.InfoBot(
		logger.InfoUpdateFulfilled,
		logger.OriginBot,
		logger.FromMenu,
		user.ID,
		user.Username,
		c.Get("start").(time.Time),
	)

	switch {
	case data == "send":
		if user.FirstName != "" {
			return c.Send(models.MsgMenuGreeting1+user.FirstName+models.MsgMenuGreeting2+"\n\n"+models.MsgMenuDescrtiption, models.ReplyMenu)
		}
		return c.Send(models.MsgMenuGreeting+"\n\n"+models.MsgMenuDescrtiption, models.ReplyMenu)

	case msgText != "":
		return c.EditOrSend(msgText+"\n\n"+models.MsgMenuDescrtiption, models.ReplyMenu)

	default:
		if user.FirstName != "" {
			return c.EditOrSend(models.MsgMenuGreeting1+user.FirstName+models.MsgMenuGreeting2+"\n\n"+models.MsgMenuDescrtiption, models.ReplyMenu)
		}
		return c.EditOrSend(models.MsgMenuGreeting+"\n\n"+models.MsgMenuDescrtiption, models.ReplyMenu)
	}
}
