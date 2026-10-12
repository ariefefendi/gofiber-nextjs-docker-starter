package controllers

import (
	"github.com/gofiber/fiber/v2"

	"myApp/go-api/config"
)

// Me (API) mengembalikan user + permission yang sedang login.
// Dipanggil frontend Next.js saat aplikasi dimuat.
func Me(c *fiber.Ctx) error {
	return c.JSON(config.AuthUser(c))
}
