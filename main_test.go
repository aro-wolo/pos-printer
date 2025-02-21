package main

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		amount   float64
		expected string
	}{
		{1234.56, "1,234.56"},
		{1000000.00, "1,000,000.00"},
		{987.65, "987.65"},
	}

	for _, test := range tests {
		result := formatMoney(test.amount)
		if result != test.expected {
			t.Errorf("formatMoney(%f) = %s; expected %s", test.amount, result, test.expected)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"12345", "54321"},
		{"hello", "olleh"},
		{"", ""},
	}

	for _, test := range tests {
		result := reverse(test.input)
		if result != test.expected {
			t.Errorf("reverse(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestGetToken(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}

	// Set environment variables
	os.Setenv("API_USERNAME", os.Getenv("API_USERNAME"))
	os.Setenv("API_KEY", os.Getenv("API_KEY"))
	os.Setenv("TOKEN_URL", os.Getenv("TOKEN_URL"))
	os.Setenv("HEADER_KEY", os.Getenv("HEADER_KEY"))
	os.Setenv("HEADER_VALUE", os.Getenv("HEADER_VALUE"))

	// Mock the HTTP response
	// ...mocking code...

	token, err := getToken()
	if err != nil {
		t.Errorf("getToken() error = %v", err)
	}
	if token == "" {
		t.Errorf("getToken() returned empty token")
	}

	// Test token caching
	tokenCache.Token = "cached_token"
	tokenCache.Expiry = time.Now().Add(1 * time.Hour)
	token, err = getToken()
	if err != nil {
		t.Errorf("getToken() error = %v", err)
	}
	if token != "cached_token" {
		t.Errorf("getToken() did not return cached token")
	}
}
