package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Md5Hash принимает строку и возвращает её MD5-хеш в виде шестнадцатеричной строки.
func Md5Hash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// GenerateToken генерирует случайный токен заданной длины (в байтах)
// и возвращает его в виде шестнадцатеричной строки.
// Функция использует криптографически стойкий генератор случайных чисел.
func GenerateToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
