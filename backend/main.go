// Package main adalah entry point API Go Fiber.
//
// Alur:
//  1. Koneksi database (config)
//  2. AutoMigrate model
//  3. Seed data dummy (opsional, SEED=true)
//  4. Daftarkan route (routes)
//  5. Jalankan server
//
// Environment variable: lihat ../.env (APP_NAME, APP_PORT, DB_*, SEED, CORS_ORIGIN)
package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"myApp/go-api/config"
	"myApp/go-api/models"
	"myApp/go-api/routes"
	"myApp/go-api/services"
)

func main() {

	// ---------- Database ----------
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	// Tambahkan model baru di sini: db.AutoMigrate(&models.Units{}, &models.Xxx{})
	if err := db.AutoMigrate(&models.Units{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Isi data dummy jika SEED=true di .env (aman dipanggil berulang)
	if config.Env("SEED", "false") == "true" {
		if err := services.SeedUnits(db); err != nil {
			log.Fatalf("seed units: %v", err)
		}
	}

	// ---------- Fiber ----------
	app := fiber.New(fiber.Config{
		AppName: config.Env("APP_NAME", "Go API"),

		// Semua error (404, panic, dll.) dikirim sebagai JSON.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var fe *fiber.Error
			if errors.As(err, &fe) {
				code = fe.Code
			}
			if code == fiber.StatusNotFound {
				return c.Status(code).JSON(fiber.Map{
					"result": "NOT_FOUND", "message": "Route not found",
				})
			}
			log.Printf("%s %s: %v", c.Method(), c.Path(), err)
			return c.Status(code).JSON(fiber.Map{
				"result": "ERROR", "message": "Internal server error",
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())

	// ---------- Routes ----------
	routes.Register(app, db)

	addr := ":" + config.Env("APP_PORT", "8000")
	log.Printf("listening on %s", addr)
	log.Fatal(app.Listen(addr))
}
