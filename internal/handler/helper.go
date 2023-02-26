package handler

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// GetIP - возвращает IP пользователя.
func GetIP(r *http.Request) (net.IP, error) {
	// получаем значения удаленного адреса пользователя из запроса
	remoteAddr := r.RemoteAddr
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, err
	}
	// парсим полученную строку
	ipFirst := net.ParseIP(ip)
	// проверяем значение заголовка X-Real-IP
	ip = r.Header.Get("X-Real-IP")

	// парсим строку адреса из заголовка и проверяем
	ipSecond := net.ParseIP(ip)
	if ipSecond == nil {
		// проверяем значение заголовка X-Forwarded-For и собираем новый адрес
		ips := r.Header.Get("X-Forwarded-For")
		splitIps := strings.Split(ips, ",")
		ip = splitIps[0]
		ipSecond = net.ParseIP(ip)
	}
	if ipFirst.Equal(ipSecond) {
		return ipFirst, nil
	}
	return nil, errors.New("ip is not real")
}
