package middleware

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

// CookiesMiddleware - middleware, проверяющяя в запросе cookie на наличие/подлинность. Если cookie нет, то генерируем новую и вставляем в Header.
// Алгоритм подписи - sha.256.
func CookiesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Ключ шифрования.
		key := []byte("AsdFrGtHyJhjErTy")
		//Проверка наличия cookie в запросе.
		if userCookie, err := r.Cookie("shortener"); err == nil {
			//Расшифровка значения cookie в срез байт.
			userCookieByte, err := hex.DecodeString(userCookie.Value)
			if err != nil {
				log.Printf("Cookie decoding: %v\n", err)
			}
			//Инициализируем алгоритм подписи HMAC.
			h := hmac.New(sha256.New, key)
			//Записываем в него полученое значение cookie.
			h.Write(userCookieByte)
			//Создаем подпись для проверки.
			sign := h.Sum(nil)
			//Проверяем на подлинность подписанной cookie.
			if hmac.Equal(userCookieByte, sign) {
				next.ServeHTTP(w, r)
				return
			}
		}
		//Если cookie нет или проверка на подлинность не пройдена, создаем новую cookie.
		//Генерация userID.
		userID, err := generateRandom(16)
		if err != nil {
			log.Printf("UserID generate: %v\n", err)
		}
		//Инициализируем алгоритм подписи HMAC.
		h := hmac.New(sha256.New, key)
		//Пишем в него сгенерированный userid.
		h.Write(userID)
		//Генерируем подписанную cookie.
		cook := h.Sum(nil)
		//Создаем cookie и передаем ее в ответ и запрос.
		cookie := &http.Cookie{
			Name:    "shortener",
			Value:   hex.EncodeToString(cook),
			Expires: time.Now().Add(time.Hour * 24),
		}
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
		next.ServeHTTP(w, r)
	})
}

// generateRandom - генератор случайных байт длинной size.
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
