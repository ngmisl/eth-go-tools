package main

import (
	"bufio"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	privateKeyLength = 64
	promptMessage    = "Enter your Ethereum private key (%d hexadecimal characters) or 'q' to quit: "
	emptyInputError  = "Private key cannot be empty"
	invalidLenError  = "Private key must be exactly %d characters long"
	invalidHexError  = "Private key must contain only hexadecimal characters"
	quitMessage      = "Exiting the program. Goodbye!"
)

// PrivateKey is a type-safe representation of an Ethereum private key
type PrivateKey string

// Validate checks if the private key is in the correct format
func (pk PrivateKey) Validate() error {
	if len(pk) == 0 {
		return errors.New(emptyInputError)
	}
	if len(pk) != privateKeyLength {
		return fmt.Errorf(invalidLenError, privateKeyLength)
	}
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", string(pk))
	if !match {
		return errors.New(invalidHexError)
	}
	return nil
}

// ConvertToAddress converts a private key to an Ethereum address
func ConvertToAddress(privateKey PrivateKey) (common.Address, error) {
	if err := privateKey.Validate(); err != nil {
		return common.Address{}, err
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(string(privateKey))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to convert private key: %v", err)
	}

	publicKey := ecdsaPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("failed to get public key")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

// handleUserInput processes user input and returns a PrivateKey and a bool indicating if the user wants to quit
func handleUserInput(reader *bufio.Reader) (PrivateKey, bool, error) {
	fmt.Printf(promptMessage, privateKeyLength)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", false, fmt.Errorf("failed to read input: %v", err)
	}

	input = strings.TrimSpace(input)
	if input == "q" || input == "Q" {
		return "", true, nil
	}

	return PrivateKey(input), false, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		privateKey, quit, err := handleUserInput(reader)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if quit {
			fmt.Println(quitMessage)
			break
		}

		address, err := ConvertToAddress(privateKey)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Public Address: %s\n\n", address.Hex())
	}
}
