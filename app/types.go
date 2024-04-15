package app

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Массив строк
type tableRows []*tableRow

// Строка таблицы
type tableRow struct {
	Id            Id `gorm:"primaryKey"`
	Region        string
	Responsible   string
	Verified      string
	VkUrl         string
	OkUrl         string
	TgUrl         string
	Reason        string
	CommentaryNpa string
	FullName      string
	Ogrn          string
	Status        string
	Commentary    string
}

// Индексы столбцов таблицы
type columns struct {
	region        int
	responsible   int
	verified      int
	vkUrl         int
	okUrl         int
	tgUrl         int
	reason        int
	commentaryNpa int
	fullName      int
	ogrn          int
	status        int
	commentary    int
}

// Хранилище данных
type storage struct {
	db *gorm.DB
}

// Привилегия
type Privilege uint8

// Идентификатор
type Id uint64

// Пользователь
type User struct {
	Id        Id `gorm:"primaryKey"`
	Privilege Privilege
	Name      string `gorm:"unique"`
	Password  string
}

// Промежуточное состояние для авторизации
type authMiddleware struct {
	storage *storage
}

// Сессия пользователя
type session struct {
	Id    Id `gorm:"primaryKey"`
	Token string

	// Для связи с пользователями
	UserID Id   `gorm:"unique"`
	User   User `gorm:"foreignKey:UserID"`
}

// Приложение
type App struct {
	router  *mux.Router
	storage *storage
	address string
	m       *authMiddleware
}

// Данные, передаваемые в шаблон страницы
type pageTemplate struct {
	// Заголовок таблицы
	TableHeader []string

	// Админ
	Admin bool
}
