// Package routes memetakan URL ke controller. Tidak berisi logika.
package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"myApp/go-api/config"
	"myApp/go-api/controllers"
	"myApp/go-api/middleware"
)

// Register mendaftarkan semua route JSON (/health dan /api/...).
func Register(app *fiber.App, db *gorm.DB) {

	health := func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err == nil {
			err = sqlDB.PingContext(c.UserContext())
		}
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "error", "database": err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "ok", "database": "connected"})
	}

	app.Get("/health", health)

	api := app.Group("/api")

	// CORS hanya aktif jika CORS_ORIGIN diisi (mis. http://localhost:3100).
	// Dengan rewrite Next.js (default kit ini) CORS tidak diperlukan.
	if origin := config.Env("CORS_ORIGIN", ""); origin != "" {
		api.Use(middleware.CORS(origin))
	}

	api.Get("/health", health)
	api.Get("/auth/me", controllers.Me)

	// ---------- Master Data : Units ----------
	units := controllers.NewUnitsController(db)
	g := api.Group("/master-data/units")
	g.Get("/data", units.GetDataAll)      // daftar + paginasi server side
	g.Get("/select", units.GetDataSelect) // satu data (?reference=)
	g.Post("/", units.Insert)             // tambah
	g.Put("/", units.Update)              // ubah
	g.Delete("/", units.Destroy)          // hapus (?reference=)

	// ---------- Test Algoritma ----------
	algo := api.Group("/test/algorithm")
	algo.Post("/bubble-sort", controllers.BubbleSort)
	algo.Post("/quick-sort", controllers.QuickSort)

}
