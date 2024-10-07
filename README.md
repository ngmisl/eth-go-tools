# Eth Go Tools

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/ngmisl/eth-go-tools/badge)](https://scorecard.dev/viewer/?uri=github.com/ngmisl/eth-go-tools)

![2024-09-20_09-41](https://github.com/user-attachments/assets/622c30ee-d8df-48bf-bfa1-095f9d58c463)

EthGoTools is a powerful and user-friendly terminal-based application built with Go. It provides a suite of Ethereum-related utilities, allowing developers and enthusiasts to manage their Ethereum accounts, generate keys, interact with Farcaster accounts, sign messages, and verify signatures—all from the comfort of your terminal.

## Table of Contents

- [Eth Go Tools](#eth-go-tools)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Installing EthGoTools](#installing-ethgotools)
  - [Usage](#usage)
    - [Running the Application](#running-the-application)
    - [Available Tools](#available-tools)
      - [1. Convert Private Key to Address](#1-convert-private-key-to-address)
      - [2. Generate New Private Key](#2-generate-new-private-key)
      - [3. Check Farcaster Account](#3-check-farcaster-account)
      - [4. Sign Message with Private Key](#4-sign-message-with-private-key)
      - [5. Verify Signature](#5-verify-signature)
  - [Configuration](#configuration)
    - [Setting Environment Variables](#setting-environment-variables)
    - [Notes](#notes)
  - [Security Considerations](#security-considerations)
  - [Contributing](#contributing)
    - [Steps to Contribute](#steps-to-contribute)
    - [Code of Conduct](#code-of-conduct)
  - [License](#license)
  - [Acknowledgements](#acknowledgements)

## Features

EthGoTools offers the following functionalities:

1. **Convert Private Key to Address**
   - Input an Ethereum private key in hexadecimal format to retrieve the corresponding Ethereum address.

2. **Generate New Private Key**
   - Generate a new Ethereum private key securely, along with its corresponding public address. **_Warning:_** The private key is displayed only once; ensure you store it securely.

3. **Check Farcaster Account**
   - Enter a Farcaster username to fetch and display profile information and recent casts using the Airstack API.

4. **Sign Message with Private Key**
   - Input a private key and a message to produce a cryptographic signature. Useful for signing transactions or authenticating messages.

5. **Verify Signature**
   - Verify the authenticity of a message signature by providing the message, signature, and the Ethereum address of the signer.

## Installation

### Prerequisites

- **Go (1.16 or later):** Ensure that [Go](https://golang.org/dl/) is installed on your system. You can verify your Go installation by running:

  ```bash
  go version
  ```

- **Git:** To clone the repository.

### Installing EthGoTools

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/ethgotools.git
   ```

2. **Navigate to the Project Directory**

   ```bash
   cd ethgotools
   ```

3. **Download Dependencies**

   EthGoTools uses Go modules for dependency management.

   ```bash
   go mod tidy
   ```

4. **Build the Application**

   ```bash
   go build -o ethgotools app.go
   ```

   This command compiles the Go application into an executable named `ethgotools`.

## Usage

### Running the Application

After building the application, you can run it directly from the terminal.

```bash
./ethgotools
```

Upon running, you'll be presented with an interactive menu to select the desired tool.

### Available Tools

#### 1. Convert Private Key to Address

**Description:** Converts an Ethereum private key to its corresponding public address.

**Steps:**

1. Select **"Convert Private Key to Address"** from the menu.
2. Enter your Ethereum private key in hexadecimal format.
3. The application will display the corresponding Ethereum address.

**Example:**

```
Ethereum Address: 0xYourEthereumAddressHere
Press Enter to continue...
```

#### 2. Generate New Private Key

**Description:** Generates a new Ethereum private key along with its corresponding public address.

**Steps:**

1. Select **"Generate New Private Key"** from the menu.
2. The application will display the newly generated private key and address.

**Warning:** **_Store your private key securely. Never share it with anyone._**

**Example:**

```
New Private Key: yournewprivatekeyhex
Corresponding Ethereum Address: 0xYourEthereumAddressHere

WARNING: Store this private key securely. Never share it with anyone!
Press Enter to continue...
```

#### 3. Check Farcaster Account

**Description:** Retrieves and displays information about a Farcaster account using the Airstack API.

**Prerequisite:** An `AIRSTACK_API_KEY`. Set this in your environment variables (see [Configuration](#configuration)).

**Steps:**

1. Select **"Check Farcaster Account"** from the menu.
2. Enter the Farcaster username you wish to check.
3. The application will display profile information and recent casts.

**Example:**

```
Results for Farcaster user 'username':

Profile Information:
Profile Name   : User's Profile Name
Follower Count : 150
Following Count: 100
FarScore       : 4.75

Recent Casts:
1. First recent cast text.
2. Second recent cast text.
...
Press Enter to continue...
```

#### 4. Sign Message with Private Key

**Description:** Signs a message using your Ethereum private key, producing a cryptographic signature.

**Steps:**

1. Select **"Sign Message with Private Key"** from the menu.
2. Enter your Ethereum private key in hexadecimal format.
3. Enter the message you wish to sign.
4. The application will display the signature.

**Example:**

```
Signature:
0xYourSignatureHere

Press Enter to return to menu...
```

#### 5. Verify Signature

**Description:** Verifies the authenticity of a signed message by checking the signature against the message and Ethereum address.

**Steps:**

1. Select **"Verify Signature"** from the menu.
2. Enter the original message that was signed.
3. Enter the signature in hexadecimal format.
4. Enter the Ethereum address of the signer.
5. The application will inform you whether the signature is valid.

**Example:**

```
Signature is valid.

Press Enter to return to menu...
```

_or_

```
Signature is invalid.

Press Enter to return to menu...
```

## Configuration

EthGoTools uses environment variables to manage sensitive information and API keys. Ensure you set these variables before running the application.

### Setting Environment Variables

Create a `.env` file in the project root directory with the following content:

```env
AIRSTACK_API_KEY=your_airstack_api_key_here
```

Alternatively, you can set environment variables directly in your shell.

**For Unix/Linux/macOS:**

```bash
export AIRSTACK_API_KEY=your_airstack_api_key_here
```

**For Windows (Command Prompt):**

```cmd
set AIRSTACK_API_KEY=your_airstack_api_key_here
```

### Notes

- The **Airstack API** is required only for the **"Check Farcaster Account"** feature.
- Ensure that your `.env` file is **never** committed to version control to protect your API keys and sensitive information.

## Security Considerations

EthGoTools handles sensitive data such as private keys and signatures. Follow these best practices to ensure your security:

- **Protect Your Private Keys:** Never share your private keys. Ensure they are stored securely and consider using hardware wallets for enhanced security.
- **Environment Variables:** Keep your `.env` file confidential. Use tools like Git's `.gitignore` to prevent accidental commits.
- **Use Trusted Environments:** Run EthGoTools on secure and trusted machines to prevent unauthorized access to your sensitive data.
- **Regular Updates:** Keep your dependencies and Go version updated to benefit from security patches and improvements.

## Contributing

Contributions are welcome! Whether you find a bug, have a feature request, or want to improve the documentation, your input is valuable.

### Steps to Contribute

1. **Fork the Repository**

   Click the "Fork" button at the top-right corner of this page to create a personal copy of the repository.

2. **Clone Your Fork**

   ```bash
   git clone https://github.com/yourusername/ethgotools.git
   cd ethgotools
   ```

3. **Create a New Branch**

   ```bash
   git checkout -b feature/YourFeatureName
   ```

4. **Make Your Changes**

   Implement your feature or bug fix.

5. **Commit Your Changes**

   ```bash
   git commit -m "Add feature: YourFeatureName"
   ```

6. **Push to Your Fork**

   ```bash
   git push origin feature/YourFeatureName
   ```

7. **Create a Pull Request**

   Navigate to the original repository and click on "Compare & pull request" to submit your changes.

### Code of Conduct

Please ensure that all contributions adhere to the [Code of Conduct](https://github.com/yourusername/ethgotools/blob/main/CODE_OF_CONDUCT.md). Be respectful, constructive, and considerate in all interactions.

## License

This project is licensed under the [MIT License](https://github.com/yourusername/ethgotools/blob/main/LICENSE).

## Acknowledgements

- **Go:** An open-source programming language that makes it easy to build simple, reliable, and efficient software.
- **Bubble Tea:** A Go framework for building delightful terminal user interfaces.
- **Lip Gloss:** A styled terminal output library for Go.
- **go-ethereum:** The official Go implementation of the Ethereum protocol.
- **Airstack:** For providing APIs to interact with Farcaster accounts.
- **Joho/Godotenv:** For loading environment variables from `.env` files.
