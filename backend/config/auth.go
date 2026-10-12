package config

import "github.com/gofiber/fiber/v2"

// AuthUser menghasilkan data user + permission untuk frontend
// (dipanggil Next.js lewat GET /api/auth/me).
//
// TODO: ini data contoh (semua akses diizinkan). Ganti dengan user dari
// session/JWT dan role + permission dari database.
//
// PENTING: permission di frontend hanya mengatur tampilan (tombol disabled).
// Penegakan sebenarnya harus dilakukan di server (middleware per route).
func AuthUser(c *fiber.Ctx) fiber.Map {

	allowAll := fiber.Map{
		"view":   true,
		"read":   true,
		"create": true,
		"update": true,
		"delete": true,
	}

	return fiber.Map{
		"user_id":   1,
		"role_id":   1,
		"role_name": "admin",
		"permissions": fiber.Map{
			// key = nama module (segmen URL, "-" menjadi "_")
			"units": allowAll,
		},
	}
}
