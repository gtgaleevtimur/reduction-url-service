// Package app аккумулирует все компоненты сервиса и запускает его работу.
// Для создании HTTPS соединения, при условии его работы
// В пакете реализована связь с  private key (.key) и public key(.crt) основанном на private key (.key).
// openssl genrsa -out server.key 2048.
// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650.
package app
