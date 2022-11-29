package middleware

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

// CookiesMiddleware - middleware, проверяющяя в запросе cookie на наличие/подлинность. Если cookie нет, то генерируем новую и вставляем в Header.
// Алгоритм подписи - AES.
func CookiesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Инициализируем буфер-переменную куда будем кодировать-декодировать cookie.
		cookieBuf := make([]byte, aes.BlockSize)
		// Секретный ключ.
		key := []byte("HdUeLk85Gp0i7pLh")
		// Проверочное слово.
		nonce := []byte("cookie")
		// Интерфейс шифрования.
		aesBlock, err := aes.NewCipher(key)
		if err != nil {
			log.Fatal(err)
		}
		// Проверка наличия cookie.
		if userCookie, err := r.Cookie("shortener"); err == nil {
			// Дешифровка cookie в строку.
			requestUserIDByte, err := hex.DecodeString(userCookie.Value)
			if err != nil {
				log.Printf("Cookie decoding: %v\n", err)
			}
			// Дешифровка строки интерфейсом шифрования.
			aesBlock.Decrypt(cookieBuf, requestUserIDByte)
			// Проверка на подлинность.
			if string(cookieBuf[len(cookieBuf)-len(nonce):]) == string(nonce) {
				next.ServeHTTP(w, r)
				return
			}
		}
		// Если cookie не обнаружено или она не прошла проверку подлиности, создаем новую cookie.
		userID, err := generateRandom(10)
		if err != nil {
			log.Printf("UserID generate: %v\n", err)
		}
		aesBlock.Encrypt(cookieBuf, append(userID, nonce...))
		cookie := &http.Cookie{
			Name: "shortener", Value: hex.EncodeToString(cookieBuf),
			Expires: time.Now().Add(time.Hour * 24),
			Path:    `/`,
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
