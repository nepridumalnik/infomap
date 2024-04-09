package app

// Ограничения для размера формы
const limitation = (32 << 20)

// Строки
const (
	// Путь до файла с логами ошибок
	errLogFile = "error_log.txt"

	// Файл СУБД
	sqliteDbFile = "data.db"

	// Соль
	salt = "qweASD123zxcASDqwe=-124-980r7102370129edqwu98"
)

// Авторизация
const (
	authorizationKey = "Authorization-Token"
	allowedPath      = "/auth"
	allowedPathHtml  = "./ui/html" + allowedPath + ".html"
)

// Столбцы таблицы
const (
	mainList      = "Томская область"
	region        = "Регион"
	responsible   = "Назначен ответственный"
	verified      = "Страница подтверждена"
	vkUrl         = "Ссылка на официальную страницу Вконтакте"
	okUrl         = "Ссылка на официальную страницу Одноклассники"
	tgUrl         = "Ссылка на официальную страницу Telegram"
	reason        = "Официальная страница не ведется на основании"
	commentaryNpa = "Комментарий по НПА"
	fullName      = "Полное наименование объекта"
	ogrn          = "ОГРН"
	status        = "Статус"
	commentary    = "Комментарий"
)

// Типы пользователей
const (
	PrivilegeUnauthorized Privilege = iota
	PrivilegeCommonUser
	PrivilegeAdmin
)
