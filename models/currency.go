package models

// OXRResponse represents the response from OpenExchangeRates API
type OXRResponse struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// CurrencySymbols represents the currency symbols response from OpenExchangeRates
type CurrencySymbols map[string]string

// ConversionRequest represents the conversion request
type ConversionRequest struct {
	FromCurrency string  `json:"from_currency" validate:"required"`
	ToCurrency   string  `json:"to_currency" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
}

// ConversionResponse represents the conversion response
type ConversionResponse struct {
	Success      bool    `json:"success"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
	Result       float64 `json:"result"`
	Rate         float64 `json:"rate"`
	Timestamp    int64   `json:"timestamp"`
}

// CachedRates represents cached currency rates (simplified - no timestamp metadata needed)
type CachedRates struct {
	OXRResponse
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
