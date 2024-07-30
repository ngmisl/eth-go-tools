package main

import (
	"crypto/ecdsa"
	"os"
	"regexp"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func TestGeneratePrivateKey(t *testing.T) {
	// Generate a new private key
	privateKey, err := generatePrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Test that the private key is not nil
	if privateKey == nil {
		t.Fatal("Generated private key is nil")
	}

	// Convert the private key to hex
	privateKeyHex := privateKeyToHex(privateKey)

	// Test that the hex representation is the correct length
	expectedHexLength := 64
	if len(privateKeyHex) != expectedHexLength {
		t.Fatalf("Private key hex length is incorrect. Expected %d, got %d", expectedHexLength, len(privateKeyHex))
	}

	// Test that the hex representation contains only valid hexadecimal characters
	validHexRegex := "^[0-9a-fA-F]+$"
	match, _ := regexp.MatchString(validHexRegex, privateKeyHex)
	if !match {
		t.Fatal("Private key hex contains invalid characters")
	}

	// Derive the public address from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal("Failed to get public key from private key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Test that the address is valid
	if !common.IsHexAddress(address.Hex()) {
		t.Fatal("Generated address is not a valid Ethereum address")
	}

	// Test that ConvertToAddress function returns the same address
	convertedAddress, err := ConvertToAddress(privateKey)
	if err != nil {
		t.Fatalf("Failed to convert private key to address: %v", err)
	}
	if convertedAddress != address {
		t.Fatalf("ConvertToAddress returned different address. Expected %s, got %s", address.Hex(), convertedAddress.Hex())
	}
}
