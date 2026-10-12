// Package controllers berisi handler HTTP: baca request, panggil service,
// kirim respons JSON. Tidak berisi query database langsung.
package controllers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"

	"myApp/go-api/services"
)

// handleServiceError memetakan error service ke respons HTTP.
func handleServiceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"result": "NOT_FOUND", "message": "Data not found",
		})
	case errors.Is(err, services.ErrDuplicateCode):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"result": "VALIDATION_ERROR",
			"errors": fiber.Map{"code": "The code has already been taken."},
		})
	default:
		return serverError(c, err)
	}
}

// validationError mengirim 422. errs boleh string atau map[field]pesan.
func validationError(c *fiber.Ctx, errs any) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"result": "VALIDATION_ERROR", "errors": errs,
	})
}

// serverError mencatat error asli di log, tetapi tidak membocorkannya ke client.
func serverError(c *fiber.Ctx, err error) error {
	log.Printf("%s %s: %v", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"result": "ERROR", "message": "Internal server error",
	})
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
