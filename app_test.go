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

	// Check that the hex representation is the correct length
	expectedHexLength := 64
	if len(privateKeyHex) != expectedHexLength {
		t.Fatalf("Private key hex length is incorrect. Expected %d, got %d", expectedHexLength, len(privateKeyHex))
	}

	// Derive the Ethereum address from the private key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Check that the address is valid
	if !common.IsHexAddress(address.Hex()) {
		t.Fatalf("Generated address is not valid: %s", address.Hex())
	}
}

func TestSignMessage(t *testing.T) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Read the private key from environment variables
	privateKeyHex := os.Getenv("TEST_PRIVATE_KEY")
	if privateKeyHex == "" {
		t.Fatal("TEST_PRIVATE_KEY not set in .env file")
	}

	message := "This is a test message"

	// Use the exported SignMessage function from app.go
	signature, err := SignMessage(privateKeyHex, message)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Check that the signature is not empty
	if signature == "" {
		t.Fatal("Generated signature is empty")
	}
}

func TestVerifySignature(t *testing.T) {
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

	message := "This is a test message"

	// Use the exported SignMessage function to create a signature
	signature, err := SignMessage(privateKeyHex, message)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Use the exported VerifySignature function to verify the signature
	valid, err := VerifySignature(message, signature, expectedAddress)
	if err != nil {
		t.Fatalf("Failed to verify signature: %v", err)
	}

	if !valid {
		t.Fatalf("Signature verification failed")
	}
}
