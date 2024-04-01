package app

import (
	"encoding/json"
	"net/http"
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
)

// Создать Bearer токен (да, странный способ)
func makeBearer(data string) string {
	token := map[string]string{
		"type": "Bearer",
		"data": data,
	}

	jsonData, _ := json.Marshal(token)
	return string(jsonData)
}

// Пока тестовая реализация проверки, чтобы было перед глазами как правильно создавать cookie
func (m *middleware) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, err := r.Cookie(authorizationKey)

		if err != nil || value == nil {
			cookie := &http.Cookie{
				Name:  authorizationKey,
				Value: makeBearer("test_token"),
			}

			http.SetCookie(w, cookie)
		}

		next.ServeHTTP(w, r)
	})
}
