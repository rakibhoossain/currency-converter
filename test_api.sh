#!/bin/bash

# Currency Converter API Test Script

BASE_URL="http://localhost:3000/api"
AUTH_TOKEN="currency_api_secure_token_2024"

echo "🚀 Testing Currency Converter API"
echo "=================================="

# Test 1: Health Check
echo "1. Health Check:"
curl -s -H "Authorization: Bearer $AUTH_TOKEN" "$BASE_URL/health" | jq '.'
echo ""

# Test 2: Get Currency Symbols
echo "2. Currency Symbols (showing first 5):"
curl -s -H "Authorization: Bearer $AUTH_TOKEN" "$BASE_URL/currencies" | jq '.symbols | to_entries | .[0:5] | from_entries'
echo ""

# Test 3: Get Currency Rates
echo "3. Currency Rates:"
curl -s -H "Authorization: Bearer $AUTH_TOKEN" "$BASE_URL/rates" | jq '.'
echo ""

# Test 4: Currency Conversion - USD to EUR
echo "4. Convert 100 USD to EUR:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USD","to_currency":"EUR","amount":100}' | jq '.'
echo ""

# Test 5: Currency Conversion - EUR to GBP
echo "5. Convert 50 EUR to GBP:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"EUR","to_currency":"GBP","amount":50}' | jq '.'
echo ""

# Test 6: Currency Conversion - JPY to USD
echo "6. Convert 1000 JPY to USD:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"JPY","to_currency":"USD","amount":1000}' | jq '.'
echo ""

# Test 7: Error Handling - Invalid Currency
echo "7. Error Test - Invalid Currency:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"INVALID","to_currency":"EUR","amount":100}' | jq '.'
echo ""

# Test 8: Error Handling - Missing Fields
echo "8. Error Test - Missing Fields:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USD","amount":100}' | jq '.'
echo ""

# Test 9: Error Handling - Missing Authorization
echo "9. Error Test - Missing Authorization:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USD","to_currency":"EUR","amount":100}' | jq '.'
echo ""

# Test 10: Error Handling - Invalid Token
echo "10. Error Test - Invalid Token:"
curl -s -X POST "$BASE_URL/convert" \
  -H "Authorization: Bearer invalid_token" \
  -H "Content-Type: application/json" \
  -d '{"from_currency":"USD","to_currency":"EUR","amount":100}' | jq '.'
echo ""

echo "✅ API Testing Complete!"
