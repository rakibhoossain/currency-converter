package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"currency-converter/models"
)

const (
	DataDir                 = "data"
	RatesCacheFile          = "rates.json"
	CurrenciesCacheFile     = "currencies.json"
	RatesCacheDuration      = 1 * time.Hour  // Update rates every hour
	CurrenciesCacheDuration = 24 * time.Hour // Update currencies daily
)

type CurrencyService struct {
	appID   string
	baseURL string
}

func NewCurrencyService(appID, baseURL string) *CurrencyService {
	// Ensure data directory exists
	os.MkdirAll(DataDir, 0755)

	service := &CurrencyService{
		appID:   appID,
		baseURL: baseURL,
	}

	// Start background cache update routines
	go service.startCacheUpdateRoutines()

	return service
}

// GetCurrencySymbols returns all currency symbols, using cache if available and valid
func (cs *CurrencyService) GetCurrencySymbols() (*models.CurrencySymbols, error) {
	// Try to get from cache first
	if symbols, err := cs.getCachedCurrencies(); err == nil && cs.isCurrenciesCacheValid() {
		return symbols, nil
	}

	// Cache is invalid or doesn't exist, fetch from API
	symbols, err := cs.fetchCurrenciesFromAPI()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currencies from API: %w", err)
	}

	// Cache the new symbols
	if err := cs.cacheCurrencies(symbols); err != nil {
		fmt.Printf("Warning: failed to cache currencies: %v\n", err)
	}

	return symbols, nil
}

// GetCurrencyRates returns currency rates, using cache if available and valid
func (cs *CurrencyService) GetCurrencyRates() (*models.OXRResponse, error) {
	// Try to get from cache first
	if rates, err := cs.getCachedRates(); err == nil && cs.isRatesCacheValid() {
		return rates, nil
	}

	// Cache is invalid or doesn't exist, fetch from API
	rates, err := cs.fetchRatesFromAPI()
	if err != nil {
		return nil, err
	}

	// Cache the new rates
	if err := cs.cacheRates(rates); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to cache rates: %v\n", err)
	}

	return rates, nil
}

// ConvertCurrency converts amount from one currency to another
func (cs *CurrencyService) ConvertCurrency(req *models.ConversionRequest) (*models.ConversionResponse, error) {
	rates, err := cs.GetCurrencyRates()
	if err != nil {
		return nil, fmt.Errorf("failed to get currency rates: %w", err)
	}

	// Validate currencies exist
	fromRate, fromExists := rates.Rates[req.FromCurrency]
	toRate, toExists := rates.Rates[req.ToCurrency]

	if !fromExists {
		return nil, fmt.Errorf("currency %s not found", req.FromCurrency)
	}
	if !toExists {
		return nil, fmt.Errorf("currency %s not found", req.ToCurrency)
	}

	// Convert: amount * (toRate / fromRate)
	// Since rates are relative to USD, we need to convert through USD
	var result float64
	var rate float64

	if req.FromCurrency == "USD" {
		result = req.Amount * toRate
		rate = toRate
	} else if req.ToCurrency == "USD" {
		result = req.Amount / fromRate
		rate = 1 / fromRate
	} else {
		// Convert from -> USD -> to
		usdAmount := req.Amount / fromRate
		result = usdAmount * toRate
		rate = toRate / fromRate
	}

	return &models.ConversionResponse{
		Success:      true,
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Amount:       req.Amount,
		Result:       result,
		Rate:         rate,
		Timestamp:    rates.Timestamp,
	}, nil
}

// fetchRatesFromAPI fetches rates from OpenExchangeRates API
func (cs *CurrencyService) fetchRatesFromAPI() (*models.OXRResponse, error) {
	url := fmt.Sprintf("%s/latest.json?app_id=%s&base=USD", cs.baseURL, cs.appID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rates models.OXRResponse
	if err := json.Unmarshal(body, &rates); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return &rates, nil
}

// getCachedRates reads rates from cache file
func (cs *CurrencyService) getCachedRates() (*models.OXRResponse, error) {
	cacheFile := filepath.Join(DataDir, RatesCacheFile)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cached models.OXRResponse
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return &cached, nil
}

// cacheRates saves rates to cache file
func (cs *CurrencyService) cacheRates(rates *models.OXRResponse) error {
	data, err := json.MarshalIndent(rates, "", "  ")
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(DataDir, RatesCacheFile)
	return os.WriteFile(cacheFile, data, 0644)
}

// isRatesCacheValid checks if cached rates are still valid based on file modification time
func (cs *CurrencyService) isRatesCacheValid() bool {
	cacheFile := filepath.Join(DataDir, RatesCacheFile)

	fileInfo, err := os.Stat(cacheFile)
	if err != nil {
		return false
	}

	age := time.Now().Sub(fileInfo.ModTime())
	return age < RatesCacheDuration
}

// fetchCurrenciesFromAPI fetches currencies from OpenExchangeRates API
func (cs *CurrencyService) fetchCurrenciesFromAPI() (*models.CurrencySymbols, error) {
	url := fmt.Sprintf("%s/currencies.json", cs.baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currencies from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var currencies models.CurrencySymbols
	if err := json.Unmarshal(body, &currencies); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return &currencies, nil
}

// getCachedCurrencies reads currencies from cache file
func (cs *CurrencyService) getCachedCurrencies() (*models.CurrencySymbols, error) {
	cacheFile := filepath.Join(DataDir, CurrenciesCacheFile)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	var cached models.CurrencySymbols
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	return &cached, nil
}

// cacheCurrencies saves currencies to cache file
func (cs *CurrencyService) cacheCurrencies(currencies *models.CurrencySymbols) error {
	data, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return err
	}

	cacheFile := filepath.Join(DataDir, CurrenciesCacheFile)
	return os.WriteFile(cacheFile, data, 0644)
}

// isCurrenciesCacheValid checks if cached currencies are still valid
func (cs *CurrencyService) isCurrenciesCacheValid() bool {
	cacheFile := filepath.Join(DataDir, CurrenciesCacheFile)

	fileInfo, err := os.Stat(cacheFile)
	if err != nil {
		return false
	}

	age := time.Now().Sub(fileInfo.ModTime())
	return age < CurrenciesCacheDuration
}

// startCacheUpdateRoutines starts background goroutines to update cache periodically
func (cs *CurrencyService) startCacheUpdateRoutines() {
	// Update rates every hour
	ratesTicker := time.NewTicker(RatesCacheDuration)
	go func() {
		for range ratesTicker.C {
			if rates, err := cs.fetchRatesFromAPI(); err == nil {
				cs.cacheRates(rates)
				fmt.Printf("Background: Updated rates cache at %v\n", time.Now().Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("Background: Failed to update rates cache: %v\n", err)
			}
		}
	}()

	// Update currencies daily
	currenciesTicker := time.NewTicker(CurrenciesCacheDuration)
	go func() {
		for range currenciesTicker.C {
			if currencies, err := cs.fetchCurrenciesFromAPI(); err == nil {
				cs.cacheCurrencies(currencies)
				fmt.Printf("Background: Updated currencies cache at %v\n", time.Now().Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("Background: Failed to update currencies cache: %v\n", err)
			}
		}
	}()

	// Initial cache population if files don't exist
	go func() {
		// Check and populate rates cache
		if !cs.isRatesCacheValid() {
			if rates, err := cs.fetchRatesFromAPI(); err == nil {
				cs.cacheRates(rates)
				fmt.Printf("Initial: Populated rates cache at %v\n", time.Now().Format("2006-01-02 15:04:05"))
			}
		}

		// Check and populate currencies cache
		if !cs.isCurrenciesCacheValid() {
			if currencies, err := cs.fetchCurrenciesFromAPI(); err == nil {
				cs.cacheCurrencies(currencies)
				fmt.Printf("Initial: Populated currencies cache at %v\n", time.Now().Format("2006-01-02 15:04:05"))
			}
		}
	}()
}
