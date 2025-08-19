package handlers

import (
	"currency-converter/models"
	"currency-converter/services"

	"github.com/gofiber/fiber/v2"
)

type CurrencyHandler struct {
	service *services.CurrencyService
}

func NewCurrencyHandler(service *services.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		service: service,
	}
}

// GetCurrencySymbols handles GET /api/currencies
func (h *CurrencyHandler) GetCurrencySymbols(c *fiber.Ctx) error {
	symbols, err := h.service.GetCurrencySymbols()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Return in the expected format with success field
	return c.JSON(fiber.Map{
		"success": true,
		"symbols": symbols,
	})
}

// GetCurrencyRates handles GET /api/rates
func (h *CurrencyHandler) GetCurrencyRates(c *fiber.Ctx) error {
	rates, err := h.service.GetCurrencyRates()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(rates)
}

// ConvertCurrency handles POST /api/convert
func (h *CurrencyHandler) ConvertCurrency(c *fiber.Ctx) error {
	var req models.ConversionRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Basic validation
	if req.FromCurrency == "" || req.ToCurrency == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Error:   "from_currency, to_currency are required and amount must be greater than 0",
		})
	}

	result, err := h.service.ConvertCurrency(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(result)
}
