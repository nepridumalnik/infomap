package app

import (
	// "encoding/json"
	"net/http"
	"strings"
)

type middleware struct {
	storage *storage
}

type session struct {
	Id    Id `gorm:"primaryKey"`
	Token string
}

const (
	authorizationKey = "Authorization-Token"
	allowedPath      = "/auth.html"
)

// Создать Bearer токен (да, странный способ)
// func makeBearer(data string) string {
// 	token := map[string]string{
// 		"type": "Bearer",
// 		"data": data,
// 	}

// 	jsonData, _ := json.Marshal(token)
// 	return string(jsonData)
// }

// Пока тестовая реализация проверки, чтобы было перед глазами как правильно создавать cookie
func (m *middleware) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.safePath(&r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Если путь не безопасный - совершаем проверку, есть ли права доступа
		value, err := r.Cookie(authorizationKey)

		if err != nil || value == nil {

			// cookie := &http.Cookie{
			// 	Name:  authorizationKey,
			// 	Value: makeBearer("test_token"),
			// }

			// http.SetCookie(w, cookie)

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
