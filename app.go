package main

import (
	"bufio"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
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

// String returns the underlying string value of the PrivateKey
func (pk PrivateKey) String() string {
	return string(pk)
}

// Validate checks if the private key is in the correct format
func (pk PrivateKey) Validate() error {
	if len(pk) == 0 {
		return errors.New(emptyInputError)
	}
	if len(pk) != privateKeyLength {
		return fmt.Errorf(invalidLenError, privateKeyLength)
	}
	match, _ := regexp.MatchString("^[0-9a-fA-F]+$", pk.String())
	if !match {
		return errors.New(invalidHexError)
	}
	return nil
}

// ConvertToECDSA converts a private key to an ECDSA private key
func ConvertToECDSA(privateKey PrivateKey) (*ecdsa.PrivateKey, error) {
	if err := privateKey.Validate(); err != nil {
		return nil, err
	}
	return crypto.HexToECDSA(privateKey.String())
}

// ConvertToAddress converts an ECDSA private key to an Ethereum address
func ConvertToAddress(ecdsaPrivateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := ecdsaPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("failed to get public key")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

// handleUserInput processes user input and returns a PrivateKey and a bool indicating if the user wants to quit
func handleUserInput(reader *bufio.Reader) (PrivateKey, bool, error) {
	fmt.Printf(promptMessage, privateKeyLength)
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", true, nil
		}
		return "", false, fmt.Errorf("failed to read input: %v", err)
	}

	input = strings.TrimSpace(input)
	if input == "q" || input == "Q" {
		return "", true, nil
	}

	return PrivateKey(input), false, nil
}

func mainMenu() {
	fmt.Println("Ethereum Toolset")
	fmt.Println("----------------")
	fmt.Println("1. Private Key Converter")
	fmt.Println("2. (Coming soon...)")
	fmt.Println("3. Quit")
}

func privateKeyConverter(reader *bufio.Reader) {
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

		ecdsaPrivateKey, err := ConvertToECDSA(privateKey)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		address, err := ConvertToAddress(ecdsaPrivateKey)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Public Address: %s\n\n", address.Hex())
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		mainMenu()
		fmt.Print("Choose an option: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		switch input {
		case "1":
			privateKeyConverter(reader)
		case "3":
			fmt.Println(quitMessage)
			return
		default:
			fmt.Println("Invalid option. Please choose again.")
		}
	}
}
