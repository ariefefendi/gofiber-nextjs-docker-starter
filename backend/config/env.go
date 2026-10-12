// Package config berisi konfigurasi aplikasi: env, database, dan auth.
package config

import "os"

// Env membaca environment variable; jika kosong dipakai fallback.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
