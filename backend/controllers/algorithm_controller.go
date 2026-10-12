package controllers

import (
	"github.com/gofiber/fiber/v2"

	"myApp/go-api/services"
)

type dataRequest struct {
	Data []int `json:"data"`
}

// POST /api/test/algorithm/bubble-sort
// body: {"data":[5,3,8,1]}
// testing : 
// Invoke-RestMethod -Method Post -Uri "http://localhost:8000/api/test/algorithm/bubble-sort" `
//   -ContentType "application/json" -Body '{"data":[64,34,25,12,22]}'

func BubbleSort(c *fiber.Ctx) error {
	var req dataRequest
	if err := c.BodyParser(&req); err != nil {
		return validationError(c, "Body JSON tidak valid.")
	}
	return c.JSON(fiber.Map{
		"result": "OK",
		"input":  "",
		"output": services.BubbleSort(req.Data),
	})
}

// POST /api/test/algorithm/quick-sort   body: {"data":[5,3,8,1]}
func QuickSort(c *fiber.Ctx) error {
	var req dataRequest
	if err := c.BodyParser(&req); err != nil {
		return validationError(c, "Body JSON tidak valid.")
	}
	return c.JSON(fiber.Map{
		"result": "OK",
		"input":  "",
		"output": services.QuickSort(req.Data),
	})
}