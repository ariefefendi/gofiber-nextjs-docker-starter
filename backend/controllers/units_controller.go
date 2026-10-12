package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"myApp/go-api/resources"
	"myApp/go-api/services"
)

type UnitsController struct {
	service *services.UnitsService
}

func NewUnitsController(db *gorm.DB) *UnitsController {
	return &UnitsController{service: services.NewUnitsService(db)}
}

// unitRequest adalah body JSON untuk Insert dan Update.
type unitRequest struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// validate memeriksa input (padanan binding:"required,max=..." di Gin).
func (r *unitRequest) validate() map[string]string {
	errs := map[string]string{}

	r.Code = strings.TrimSpace(r.Code)
	r.Name = strings.TrimSpace(r.Name)

	switch {
	case r.Code == "":
		errs["code"] = "The code field is required."
	case len(r.Code) > 50:
		errs["code"] = "The code may not be greater than 50 characters."
	}

	switch {
	case r.Name == "":
		errs["name"] = "The name field is required."
	case len(r.Name) > 100:
		errs["name"] = "The name may not be greater than 100 characters."
	}

	return errs
}

func (r unitRequest) input() services.UnitInput {
	return services.UnitInput{
		Code:        r.Code,
		Name:        r.Name,
		Description: r.Description,
		IsActive:    r.IsActive,
	}
}

// GetDataAll (API) daftar data dengan paginasi server side.
// Query: filtertext, start (offset), length (jumlah), draw.
//
// Respons: { draw, total, filtered, data: [...] }
func (h *UnitsController) GetDataAll(c *fiber.Ctx) error {

	filter := strings.TrimSpace(c.Query("filtertext"))
	start := clamp(c.QueryInt("start", 0), 0, 1_000_000)
	length := clamp(c.QueryInt("length", 10), 1, 100)
	draw := c.QueryInt("draw", 0)

	items, total, filtered, err := h.service.Paginate(filter, start, length)
	if err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{
		"draw":     draw,
		"total":    total,
		"filtered": filtered,
		"data":     resources.NewUnitCollection(items),
	})
}

// GetDataSelect (API) mengambil satu data. Query: reference.
func (h *UnitsController) GetDataSelect(c *fiber.Ctx) error {

	unit, err := h.service.FindByReference(c.Query("reference"))
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(resources.NewUnitResource(*unit))
}

// Insert (API) menambah data baru. Body: unitRequest.
func (h *UnitsController) Insert(c *fiber.Ctx) error {

	var req unitRequest
	if err := c.BodyParser(&req); err != nil {
		return validationError(c, "invalid JSON body")
	}
	if errs := req.validate(); len(errs) > 0 {
		return validationError(c, errs)
	}

	if _, err := h.service.Create(req.input()); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{"result": "OK"})
}

// Update (API) mengubah data. Body: unitRequest (wajib "id").
func (h *UnitsController) Update(c *fiber.Ctx) error {

	var req unitRequest
	if err := c.BodyParser(&req); err != nil {
		return validationError(c, "invalid JSON body")
	}
	if req.ID == "" {
		return validationError(c, "id is required")
	}
	if errs := req.validate(); len(errs) > 0 {
		return validationError(c, errs)
	}

	if err := h.service.Update(req.ID, req.input()); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{"result": "OK"})
}

// Destroy (API) menghapus data. Query: reference.
func (h *UnitsController) Destroy(c *fiber.Ctx) error {

	if err := h.service.Delete(c.Query("reference")); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{"result": "OK"})
}
