package app

import (
	"encoding/json"
	"net/http"
	"strings"
)

type middleware struct {
	storage *storage
}

type session struct {
	Id    Id `gorm:"primaryKey"`
	Token string

	// Для связи с пользователями
	UserID Id   `gorm:"unique"`
	User   User `gorm:"foreignKey:UserID"`
}

const (
	authorizationKey = "Authorization-Token"
	allowedPath      = "/auth"
	allowedPathHtml  = "./ui/html" + allowedPath + ".html"
)

// Создать Bearer токен (да, странный способ)
func makeBearer(data string) string {
	token := map[string]string{
		"type": "Bearer",
		"data": data,
	}

	jsonData, _ := json.Marshal(token)
	return Sha512(string(jsonData))
}

// Обработчик авторизации
func (m *middleware) authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, allowedPathHtml)

	case "POST":
		name := r.FormValue("user")
		password := r.FormValue("password")

		data := &User{Name: name, Password: Sha512(password)}
		var user User

		result := m.storage.db.Where(data, "name", "password").Find(&user)

		if result.RowsAffected == 1 {
			session := session{Token: makeBearer(password), UserID: user.Id}
			result := m.storage.db.Where(&session).First(&session)

			if result.RowsAffected == 0 {
				m.storage.db.Create(&session)
			} else {
				m.storage.db.Save(&session)
			}

			cookie := &http.Cookie{
				Name:  authorizationKey,
				Value: makeBearer(password),
			}

			http.SetCookie(w, cookie)

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		http.Error(w, "wrong data", http.StatusUnauthorized)
	}
}

func (m *middleware) unauthHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(authorizationKey)

	if err != nil || token.Value == "" {
		return
	}

	cookie := &http.Cookie{
		Name: authorizationKey,
	}

	m.storage.db.Delete(&session{Token: token.Value}, "token")
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Пока тестовая реализация проверки, чтобы было перед глазами как правильно создавать cookie
func (m *middleware) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.safePath(&r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Если путь не безопасный - совершаем проверку, есть ли права доступа
		token, err := r.Cookie(authorizationKey)

		if err != nil || token == nil {
			// При неудачной проверке перенаправляем на страницу авторизации
			http.Redirect(w, r, allowedPath, http.StatusSeeOther)
			return
		}

		s := session{Token: token.Value}
		result := m.storage.db.Where(&s).First(&s)

		if result.RowsAffected == 0 || s.Token != token.Value {
			// При неудачной проверке перенаправляем на страницу авторизации
			http.Redirect(w, r, allowedPath, http.StatusSeeOther)
			return
		}

		// При успешной проверке прав позволяем совершить запрос
		next.ServeHTTP(w, r)
	})
}

// Проверка допустимых путей
func (m *middleware) safePath(url *string) bool {
	return *url == allowedPath || strings.HasSuffix(*url, ".css") || strings.HasSuffix(*url, ".js")
}
