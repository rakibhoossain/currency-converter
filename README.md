# Currency Converter API

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/rakibhoossain/currency-converter)](https://github.com/rakibhoossain/currency-converter/releases)

A Go REST API built with Fiber framework that provides currency conversion services using the OpenExchangeRates API with intelligent caching and Bearer token authentication.

## Features

- **Currency Symbols**: Get all available currency symbols and their descriptions
- **Currency Rates**: Get current exchange rates (cached for 1 hour)
- **Currency Conversion**: Convert amounts between different currencies
- **Smart Caching**: Automatically caches exchange rates for 1 hour and currency symbols for 24 hours
- **Background Updates**: Automatic cache refresh in the background
- **Authentication**: Bearer token authentication for all API endpoints
- **Error Handling**: Comprehensive error handling and validation

## Authentication

All API endpoints require authentication using a Bearer token. Include the token in the Authorization header:

```
Authorization: Bearer your_auth_token_here
```

## API Endpoints

### 1. Get Currency Symbols
```
GET /api/currencies
```
Returns all available currency symbols and their descriptions.

**Response:**
```json
{
  "success": true,
  "symbols": {
    "USD": "United States Dollar",
    "EUR": "Euro",
    "GBP": "British Pound Sterling",
    ...
  }
}
```

### 2. Get Currency Rates
```
GET /api/rates
```
Returns current exchange rates (base: USD).

**Response:**
```json
{
  "disclaimer": "Usage subject to terms: https://openexchangerates.org/terms",
  "license": "https://openexchangerates.org/license",
  "timestamp": 1755630000,
  "base": "USD",
  "rates": {
    "EUR": 0.858563,
    "GBP": 0.741542,
    "JPY": 147.514333,
    ...
  }
}
```

### 3. Convert Currency
```
POST /api/convert
```

**Request Body:**
```json
{
  "from_currency": "USD",
  "to_currency": "EUR",
  "amount": 100
}
```

**Response:**
```json
{
  "success": true,
  "from_currency": "USD",
  "to_currency": "EUR",
  "amount": 100,
  "result": 85.856,
  "rate": 0.858563,
  "timestamp": 1755630000
}
```

### 4. Health Check
```
GET /api/health
```
Returns API health status.

## 🚀 Quick Start

### Option 1: Download Pre-built Binary (Recommended)

1. **Download the latest release** for your platform from [Releases](https://github.com/rakibhoossain/currency-converter/releases)

2. **Extract the archive:**
```bash
# Linux/macOS
tar -xzf currency-converter-v1.0.0-linux-amd64.tar.gz
cd currency-converter-v1.0.0-linux-amd64/

# Windows
# Extract currency-converter-v1.0.0-windows-amd64.zip
```

3. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your settings
```

4. **Run the application:**
```bash
# Linux/macOS
./start.sh

# Windows
start.bat

# Or run directly
./currency-converter-linux-amd64  # Linux
./currency-converter-darwin-amd64 # macOS
currency-converter-windows-amd64.exe # Windows
```

### Option 2: Build from Source

1. **Clone the repository:**
```bash
git clone https://github.com/rakibhoossain/currency-converter.git
cd currency-converter
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Set up environment variables:**
```bash
cp .env.example .env
```
Edit `.env` and add your OpenExchangeRates App ID and authentication token:
```
OXR_APP_ID=your_actual_app_id_here
OXR_BASE_URL=https://openexchangerates.org/api
AUTH_TOKEN=your_secure_auth_token_here
PORT=3000
```

4. **Run the application:**
```bash
go run main.go
```

5. **Test the API:**
```bash
./test_api.sh
```

The API will be available at `http://localhost:3000`

### Option 3: Build Multi-platform Binaries

```bash
./build.sh
```
This creates binaries for Linux, Windows, and macOS in the `dist/` directory.

## Docker Deployment

1. **Build the Docker image:**
```bash
docker build -t currency-converter .
```

2. **Run with Docker:**
```bash
docker run -p 3000:3000 -e OXR_APP_ID=your_app_id_here -e AUTH_TOKEN=your_auth_token_here currency-converter
```

3. **Run with Docker Compose:**
Create a `docker-compose.yml`:
```yaml
version: '3.8'
services:
  currency-converter:
    build: .
    ports:
      - "3000:3000"
    environment:
      - OXR_APP_ID=your_app_id_here
      - OXR_BASE_URL=https://openexchangerates.org/api
      - AUTH_TOKEN=your_auth_token_here
      - PORT=3000
    volumes:
      - ./data:/root/data
```

Then run:
```bash
docker-compose up
```

## 📁 Project Structure

```
currency_converter/
├── main.go                 # Application entry point
├── models/
│   └── currency.go         # Data models
├── services/
│   └── currency_service.go # Business logic with caching
├── handlers/
│   └── currency_handler.go # HTTP handlers
├── middleware/
│   └── auth.go            # Authentication middleware
├── routes/
│   └── routes.go          # Route definitions
├── data/                  # Data directory (auto-created)
│   ├── rates.json        # Cached exchange rates
│   └── currencies.json   # Cached currency symbols
├── dist/                  # Build artifacts (created by build.sh)
├── build.sh              # Multi-platform build script
├── test_api.sh           # API testing script
├── .env.example          # Environment variables template
├── go.mod               # Go module file
└── README.md           # This file
```

## Caching Strategy

- **Exchange rates**: Cached for 1 hour and updated automatically in the background
- **Currency symbols**: Cached for 24 hours and updated automatically in the background
- Cache files: `data/rates.json` and `data/currencies.json`
- Cache validity is determined by file modification time (no embedded timestamps needed)
- If cache is expired or missing, fresh data is fetched from OpenExchangeRates API
- Background goroutines automatically refresh cache at specified intervals

## Error Handling

The API provides comprehensive error handling:
- Invalid request format
- Missing required fields
- Invalid currency codes
- API failures
- Network errors

All errors return a consistent format:
```json
{
  "success": false,
  "error": "Error description"
}
```

## 🛠️ Development

### Building from Source
```bash
go build -o currency-converter .
```

### Running Tests
```bash
./test_api.sh
```

### Building for Multiple Platforms
```bash
./build.sh
```

## 📦 Releases

Pre-built binaries are available for:
- **Linux** (amd64)
- **Windows** (amd64)
- **macOS** (amd64, arm64)

Download from [GitHub Releases](https://github.com/rakibhoossain/currency-converter/releases)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📋 Requirements

- Go 1.21 or higher
- OpenExchangeRates API account (free tier available)

## 🔗 Dependencies

- [Fiber v2](https://github.com/gofiber/fiber) - Web framework
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- [OpenExchangeRates](https://openexchangerates.org/) for providing currency data
- [Fiber](https://gofiber.io/) for the excellent web framework