package middleware

import (
	"compress/gzip"
	"net/http"
)

// Decompress - middleware, распаковывающая GZIP запросы.
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Если запрос сжат gzip,то заменяем r.Reader на gzip.Reader.
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			g, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = g
			defer g.Close()
		}
		//Если запрос не сжат ,то передаем дальше.
		next.ServeHTTP(w, r)
	})
}
