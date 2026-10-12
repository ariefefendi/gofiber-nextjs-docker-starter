// Package middleware berisi fungsi yang berjalan sebelum controller.
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS mengizinkan request lintas origin ke route API.
// Hanya perlu jika browser memanggil API langsung (tanpa rewrite Next.js).
func CORS(allowOrigin string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: allowOrigin,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	})
}
