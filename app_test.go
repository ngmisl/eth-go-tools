package main

import (
	"os"
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
	privateKeyHex := os.Getenv("TEST_PRIVATE_KEY")
	expectedAddress := os.Getenv("TEST_EXPECTED_ADDRESS")

	if privateKeyHex == "" || expectedAddress == "" {
		t.Fatal("TEST_PRIVATE_KEY or TEST_EXPECTED_ADDRESS not set in .env file")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		t.Fatalf("Failed to convert private key to ECDSA: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	if address.Hex() != expectedAddress {
		t.Errorf("Address mismatch. Expected %s, got %s", expectedAddress, address.Hex())
	}
}

func TestGeneratePrivateKey(t *testing.T) {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Test that the private key is not nil
	if privateKey == nil {
		t.Fatal("Generated private key is nil")
	}

	// Convert the private key to hex
	privateKeyHex := common.Bytes2Hex(crypto.FromECDSA(privateKey))

	// Test that the hex representation is the correct length
	expectedHexLength := 64
	if len(privateKeyHex) != expectedHexLength {
		t.Fatalf("Private key hex length is incorrect. Expected %d, got %d", expectedHexLength, len(privateKeyHex))
	}

	// Derive the Ethereum address from the private key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Test that the address is valid
	if !common.IsHexAddress(address.Hex()) {
		t.Fatalf("Generated address is not valid: %s", address.Hex())
	}
}
