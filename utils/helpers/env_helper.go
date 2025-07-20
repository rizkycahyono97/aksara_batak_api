package helpers

import (
	"os"
	"strconv"
	"strings"
)

// Fungsi ini digunakan untuk membaca nilai dari sebuah variabel lingkungan yang berisi daftar nilai,
// dipisahkan oleh koma (,), dan mengembalikannya sebagai slice ([]string).
func ParseEnvList(key string) []string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return []string{}
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// Fungsi ini digunakan untuk membaca nilai boolean dari sebuah variabel lingkungan.
// Jika variabel tidak ada, nilainya kosong, atau tidak dapat diubah menjadi boolean, fungsi akan mengembalikan nilai default yang ditentukan.
func GetEnvBool(key string, defaultValue bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	parsedVal, err := strconv.ParseBool(strings.ToLower(val))
	if err != nil {
		return defaultValue
	}
	return parsedVal
}
