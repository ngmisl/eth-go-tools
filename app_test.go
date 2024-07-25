package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestPrivateKeyConverter(t *testing.T) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Read the private key and expected address from environment variables
	privateKey := os.Getenv("TEST_PRIVATE_KEY")
	expectedAddress := os.Getenv("TEST_EXPECTED_ADDRESS")

	if privateKey == "" || expectedAddress == "" {
		t.Fatal("TEST_PRIVATE_KEY or TEST_EXPECTED_ADDRESS not set in .env file")
	}

	ecdsaPrivateKey, err := ConvertToECDSA(PrivateKey(privateKey))
	if err != nil {
		t.Fatalf("Failed to convert private key to ECDSA: %v", err)
	}

	address, err := ConvertToAddress(ecdsaPrivateKey)
	if err != nil {
		t.Fatalf("Failed to convert ECDSA key to address: %v", err)
	}

	if address.Hex() != expectedAddress {
		t.Errorf("Address mismatch. Expected %s, got %s", expectedAddress, address.Hex())
	}
}
