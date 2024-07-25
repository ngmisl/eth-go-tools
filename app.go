package main

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PrivateKey string

func (pk PrivateKey) String() string {
	return string(pk)
}

func ConvertToECDSA(privateKey PrivateKey) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(privateKey.String())
}

func ConvertToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nEthereum Toolkit")
		fmt.Println("1. Convert Private Key to Address")
		fmt.Println("2. Exit")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter your private key: ")
			privateKeyInput, _ := reader.ReadString('\n')
			privateKeyInput = strings.TrimSpace(privateKeyInput)

			privateKey := PrivateKey(privateKeyInput)

			ecdsaPrivateKey, err := ConvertToECDSA(privateKey)
			if err != nil {
				fmt.Printf("Error converting private key: %v\n", err)
				continue
			}

			address, err := ConvertToAddress(ecdsaPrivateKey)
			if err != nil {
				fmt.Printf("Error converting to address: %v\n", err)
				continue
			}

			fmt.Printf("Ethereum Address: %s\n", address.Hex())

		case "2":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
