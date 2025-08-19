package routes

import (
	"currency-converter/handlers"
	"currency-converter/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, currencyHandler *handlers.CurrencyHandler, authToken string) {
	// API routes
	api := app.Group("/api")

	// Apply authentication middleware to all API routes
	api.Use(middleware.AuthMiddleware(authToken))

	// Currency routes
	api.Get("/currencies", currencyHandler.GetCurrencySymbols)
	api.Get("/rates", currencyHandler.GetCurrencyRates)
	api.Post("/convert", currencyHandler.ConvertCurrency)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Currency Converter API is running",
		})
	})
}
