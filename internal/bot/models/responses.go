package models

var (
	// Отправляется при ошибке.
	MsgOnError = "× Что-то пошло не так... Свяжись и поделись ошибкой с @xoticdsign - моим создателем!\n\nОшибка: "
)

var (
	// Стандартная первая половина сообщения из главного меню.
	MsgMenuGreeting = "< Привет, я помогу тебе сократить твои ссылки! >\n\n"

	// Вторая половина сообщения из главного меню.
	MsgMenuDescrtiption = "| Мой создатель: @xoticdsign\n\nОтправь ссылку в чат со мной, чтобы я ее сократил. Также ты можешь воспользоваться одной из команд ниже, чтобы управлять своими Shorties!\n\nПример ссылки:\nhttps://trex-runner.com/night/"
)

var (
	// Отправляется, если удалось создать сокращенную ссылку. Используется, как первая половина сообщения из главного меню.
	SuccessNew = "+ Создал для тебя Shorty! Я прикреплю ее ниже.\n\n"

	// Отправляется, если не удалось создать сокращенную ссылку. Причина: пользователь отправил ссылку, сокращенную ботом.
	FailedNewCantShortShorty = "× Я не могу сократить Shorty!\n\nПопробуй отправить правильную ссылку или вернись в меню."

	// Отправляется, если не удалось создать сокращенную ссылку. Причина: некорректная ссылка.
	FailedNewIncorrectURL = "× Прости, но это некорректная ссылка!\n\nПопробуй отправить правильную ссылку или вернись в меню."

	// Отправляется, если не удалось создать сокращенную ссылку. Причина: ссылка уже существует.
	FailedNewDuplicate = "× Такая Shorty уже существует!\n\nОтправь другую ссылку или вернись в меню."

	// Отправляется, если не удалось создать сокращенную ссылку. Причина: пользователь превысил лимит созданных ссылок (5).
	FailedNewLimitExceeded = "× Ты превысил превысил лимит возможных Shorties!\n\nУдали ненужную Shorty."
)

var (
	// Отправляется, если пользователь зашел в раздел "Мои Shorties".
	MsgListShorties = "Я нашел все твои Shorties!\n\nТы можешь выбрать одну из них, чтобы узнать больше подробностей."

	// Отправляется, если пользователь выбрал ссылку в разделе "Мои Shorties".
	MsgShortyInfo = "Вот вся доступная информация об этой Shorty!"
)

var (
	// Отправляется, если пользователь зашел в раздел "Удалить Shorty".
	MsgDeleteShorty = "Закрепил под этим сообщением все созданные тобой Shorties!\n\nВыбери ту Shorty, которая тебе больше не нужна. Не переживай, я переспрошу у тебя, точно ли ты хочешь ее удалить!"

	// Отправляется, если пользователь выбрал ссылку в разделе "Удалить Shorty".
	MsgDeleteShortyPrompt = "Ты уверен, что хочешь удалить эту Shorty? Это действие отменить не получится!"

	// Отправляется, если удалось удалить сокращенную ссылку. Используется, как первая половина сообщения из главного меню.
	SuccessDelete = "- Shorty успешно удалена!"
)

var (
	// Отправляется, если пользователь использовал неподдерживаемую команду.
	FailedGlobalUnsupportedCmd = "× Такой команды я не знаю!"

	// Отправляется, если у пользователя нет сокращенных ссылок.
	FailedGlobalNoShorties = "× У тебя нет доступных Shorties!\n\nПопробуй создать парочку. Для того, чтобы узнать, как это сделать - вернись в главное меню."

	// Отправляется, если пользователь не установил Username в настройках Телеграмма.
	FailedGlobalUsernameAbsent = "× У тебя не установлено имя пользователя, я не смогу тебя запомнить!\n\nУстанови имя пользователя в настройках Telegram."
)
