package middlewares

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

//	CookiesMiddleware - middleware, проверяющяя в запросе cookie на наличие/подлинность. Если cookie нет, то генерируем новую и вставляем в Header.
// Алгоритм шифрования - AES.
func CookiesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Срез байт в качестве буфера.
		//authCookie := make([]byte, aes.BlockSize)
		//Ключ шифрования.
		key := []byte("AsdFrGtHyJhjErTy")
		//Проверочное слово.
		//nonce := []byte("cookie")
		//Инициализация интерфейся для симметричного шифрования.
		//aesBlock, err := aes.NewCipher(key)
		//if err != nil {
		//	log.Println(err)
		//}
		//Проверка наличия cookie в запросе.
		if userCookie, err := r.Cookie("shortener"); err == nil {
			//Расшифровка значения cookie в срез байт.
			userCookieByte, err := hex.DecodeString(userCookie.Value)
			if err != nil {
				log.Printf("Cookie decoding: %v\n", err)
			}
			//Расшифровываем полученный срез байт с помощью интерфейса симметричного шифрования.z`
			//aesBlock.Decrypt(authCookie, userCookieByte)
			h := hmac.New(sha256.New, key)
			h.Write(userCookieByte)
			sign := h.Sum(nil)
			//Проверяем на подлинность сравнением конца расшированного означения от длинны проверчного слова и проверчного слова.
			//if string(authCookie[len(authCookie)-len(nonce):]) == string(nonce) {
			//next.ServeHTTP(w, r)
			//return
			//}
			if hmac.Equal(userCookieByte, sign) {
				next.ServeHTTP(w, r)
				return
			}
		}
		//Если cookie нет или проверка на подлинность не пройдена,создаем новую cookie.
		//Генерация userID.
		userID, err := generateRandom(16)
		if err != nil {
			log.Printf("UserID generate: %v\n", err)
		}
		/*
			//Шифрование userID+проверочное слово.
			aesBlock.Encrypt(authCookie, append(userID, nonce...))
			//Инициализация и вставка новой cookie.
		*/
		h := hmac.New(sha256.New, key)
		h.Write(userID)
		cook := h.Sum(nil)
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

//generateRandom - генератор случайных байт длинной size.
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
