package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Создать Bearer токен (да, странный способ)
func makeBearer(data string) string {
	token := struct {
		TokenType string
		Data      string
		Timestamp int64
	}{
		TokenType: "Bearer",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	jsonData, _ := json.Marshal(token)
	return Sha512(string(jsonData))
}

// Обработчик авторизации
func (m *authMiddleware) authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		_, err := m.isAuthorized(r)

		if err != nil {
			http.ServeFile(w, r, allowedPathHtml)
			return
		}

		http.ServeFile(w, r, allowedPathHtml)

	case "POST":
		name := r.FormValue("user")
		password := r.FormValue("password")

		data := &User{Name: name, Password: Sha512(password)}
		var user User

		result := m.storage.db.Where(data, "name", "password").Find(&user)

		if result.RowsAffected != 1 {
			http.Error(w, "wrong data", http.StatusUnauthorized)
			return
		}

		session := session{UserID: user.Id}
		result = m.storage.db.Find(&session).Where("user_id")

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusBadRequest)
		}

		session.Token = makeBearer(password)

		if result.RowsAffected == 0 {
			m.storage.db.Create(&session)
		} else {
			m.storage.db.Save(&session)
		}

		cookie := &http.Cookie{
			Name:    authorizationKey,
			Value:   session.Token,
			Expires: time.Now().AddDate(0, 0, 3),
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Разлогинивание
func (m *authMiddleware) unauthHandler(w http.ResponseWriter, r *http.Request) {
	token, err := m.isAuthorized(r)

	if err != nil || token.Value == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name: authorizationKey,
	}

	m.storage.db.Delete(&session{Token: token.Value}, "token")
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Авторизован ли пользователь
func (m *authMiddleware) isAuthorized(r *http.Request) (*http.Cookie, error) {
	token, err := r.Cookie(authorizationKey)

	if err != nil || token.Value == "" {
		return nil, err
	}

	return token, nil
}

// Пока тестовая реализация проверки, чтобы было перед глазами как правильно создавать cookie
func (m *authMiddleware) authMiddleware(next http.Handler) http.Handler {
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
func (m *authMiddleware) safePath(url *string) bool {
	return *url == allowedPath || strings.HasSuffix(*url, ".css") || strings.HasSuffix(*url, ".js")
}
