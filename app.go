// app.go

package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"example.com/ethgotools/airstack"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices  []string
	cursor   int
	selected string
	quitting bool
	content  string
	input    string
	input2   string
	input3   string
	state    string
	step     int
}

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4")).
	Padding(0, 1)

var menuStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFA500"))

var inputStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00FF7F"))

func initialModel() model {
	return model{
		choices: []string{
			"Convert Private Key to Address",
			"Generate New Private Key",
			"Check Farcaster Account",
			"Sign Message with Private Key",
			"Verify Signature",
			"Quit",
		},
		state: "menu",
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found or error loading it.")
	}

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case "menu":
		return m.updateMenu(msg)
	case "convert":
		return m.updateConvert(msg)
	case "generate":
		return m.updateGenerate()
	case "farcaster":
		return m.updateFarcaster(msg)
	case "sign":
		return m.updateSign(msg)
	case "verify":
		return m.updateVerify(msg)
	case "display":
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			m.state = "menu"
			m.content = ""
			m.input = ""
			m.input2 = ""
			m.input3 = ""
			m.step = 0
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case "menu":
		return m.viewMenu()
	case "convert":
		return m.viewConvert()
	case "generate":
		return m.viewGenerate()
	case "farcaster":
		return m.viewFarcaster()
	case "sign":
		return m.viewSign()
	case "verify":
		return m.viewVerify()
	case "display":
		return m.viewDisplay()
	default:
		return "Unknown state"
	}
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.choices[m.cursor]
			switch m.cursor {
			case 0:
				m.state = "convert"
				m.input = ""
				m.content = ""
			case 1:
				m.state = "generate"
			case 2:
				m.state = "farcaster"
				m.input = ""
				m.content = ""
			case 3:
				m.state = "sign"
				m.input = ""
				m.input2 = ""
				m.content = ""
				m.step = 0
			case 4:
				m.state = "verify"
				m.input = ""
				m.input2 = ""
				m.input3 = ""
				m.content = ""
				m.step = 0
			case 5:
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) viewMenu() string {
	s := titleStyle.Render("Ethereum Tools Menu") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, menuStyle.Render(choice))
	}

	s += "\nPress q to quit.\n"

	return s
}

func (m model) updateConvert(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.state = "menu"
		case tea.KeyEnter:
			if len(m.input) == 0 {
				m.content = "Error: Private key cannot be empty."
				return m, nil
			}

			privateKeyHex := strings.TrimSpace(m.input)
			privateKey, err := crypto.HexToECDSA(privateKeyHex)
			if err != nil {
				m.content = fmt.Sprintf("Error converting private key: %v", err)
				return m, nil
			}

			address := crypto.PubkeyToAddress(privateKey.PublicKey)
			m.content = fmt.Sprintf("Ethereum Address: %s", address.Hex())
			m.state = "display"
			return m, nil
		case tea.KeyRunes:
			m.input += string(msg.Runes)
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		}
	}
	return m, nil
}

func (m model) viewConvert() string {
	s := titleStyle.Render("Convert Private Key to Address") + "\n\n"
	s += "Enter your Ethereum private key (in hex format) or press Esc to cancel:\n"
	s += inputStyle.Render(m.input)
	if m.content != "" {
		s += "\n\n" + m.content + "\n\nPress Enter to continue..."
	}
	return s
}

func (m model) updateGenerate() (tea.Model, tea.Cmd) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		m.content = fmt.Sprintf("Error generating private key: %v", err)
	} else {
		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := fmt.Sprintf("%x", privateKeyBytes)
		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		m.content = fmt.Sprintf("New Private Key: %s\nCorresponding Ethereum Address: %s\n\nWARNING: Store this private key securely. Never share it with anyone!", privateKeyHex, address.Hex())
	}
	m.state = "display"
	return m, nil
}

func (m model) viewGenerate() string {
	s := titleStyle.Render("Generate New Private Key") + "\n\n"
	if m.content != "" {
		s += m.content + "\n\nPress Enter to continue..."
	}
	return s
}

func (m model) updateFarcaster(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.state = "menu"
		case tea.KeyEnter:
			if len(m.input) == 0 {
				m.content = "Error: Farcaster username cannot be empty."
				return m, nil
			}

			// Start fetching data
			m.content = "Waiting for answer..."
			// Store the input locally to avoid race conditions
			input := m.input
			return m, func() tea.Msg {
				// Perform API call here
				client := airstack.NewClient()
				apiKey := os.Getenv("AIRSTACK_API_KEY")
				if apiKey == "" {
					return "Error: AIRSTACK_API_KEY not set."
				}
				client.SetAPIKey(apiKey)

				fname := strings.TrimSpace(input)

				result, err := client.QueryFarcasterAccount(fname)
				if err != nil {
					return fmt.Sprintf("Error querying Airstack API: %v", err)
				}

				// Format the results
				if len(result.Data.Socials.Social) == 0 && len(result.Data.FarcasterCasts.Cast) == 0 {
					return "No data found for the provided Farcaster username."
				}

				output := formatFarcasterData(fname, result)
				return output
			}
		case tea.KeyRunes:
			m.input += string(msg.Runes)
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		}
	case string:
		m.content = msg
		m.state = "display"
		return m, nil
	}

	return m, nil
}

func (m model) viewFarcaster() string {
	s := titleStyle.Render("Check Farcaster Account") + "\n\n"
	s += "Enter Farcaster username or press Esc to cancel:\n"
	s += inputStyle.Render(m.input)
	if m.content != "" {
		s += "\n\n" + m.content + "\n\nPress Enter to continue..."
	}
	return s
}

func (m model) updateSign(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.input = ""
			m.input2 = ""
			m.content = ""
			m.step = 0
			m.state = "menu"
		case tea.KeyEnter:
			if m.step == 0 {
				if len(m.input) == 0 {
					m.content = "Error: Private key cannot be empty."
					return m, nil
				}
				// Validate private key
				_, err := crypto.HexToECDSA(strings.TrimSpace(m.input))
				if err != nil {
					m.content = "Error: Invalid private key format."
					return m, nil
				}
				m.step = 1
			} else if m.step == 1 {
				if len(m.input2) == 0 {
					m.content = "Error: Message cannot be empty."
					return m, nil
				}
				// Sign the message
				privateKeyHex := strings.TrimSpace(m.input)
				message := m.input2
				signature, err := SignMessage(privateKeyHex, message)
				if err != nil {
					m.content = fmt.Sprintf("Error signing message: %v", err)
					return m, nil
				}
				m.content = fmt.Sprintf("Signature:\n%s", signature)
				m.state = "display"
				return m, nil
			}
		case tea.KeyRunes:
			if m.step == 0 {
				m.input += string(msg.Runes)
			} else if m.step == 1 {
				m.input2 += string(msg.Runes)
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if m.step == 0 && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			} else if m.step == 1 && len(m.input2) > 0 {
				m.input2 = m.input2[:len(m.input2)-1]
			}
		}
	}
	return m, nil
}

func (m model) viewSign() string {
	s := titleStyle.Render("Sign Message with Private Key") + "\n\n"
	if m.step == 0 {
		s += "Enter your Ethereum private key (in hex format) or press Esc to cancel:\n"
		s += inputStyle.Render(m.input)
		if m.content != "" {
			s += "\n\n" + m.content + "\n\nPress Enter to continue..."
		}
	} else if m.step == 1 {
		s += "Enter the message you wish to sign or press Esc to cancel:\n"
		s += inputStyle.Render(m.input2)
		if m.content != "" {
			s += "\n\n" + m.content + "\n\nPress Enter to continue..."
		}
	}
	return s
}

func (m model) updateVerify(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.input = ""
			m.input2 = ""
			m.input3 = ""
			m.content = ""
			m.step = 0
			m.state = "menu"
		case tea.KeyEnter:
			if m.step == 0 {
				if len(m.input) == 0 {
					m.content = "Error: Message cannot be empty."
					return m, nil
				}
				m.step = 1
			} else if m.step == 1 {
				if len(m.input2) == 0 {
					m.content = "Error: Signature cannot be empty."
					return m, nil
				}
				m.step = 2
			} else if m.step == 2 {
				if len(m.input3) == 0 {
					m.content = "Error: Ethereum address cannot be empty."
					return m, nil
				}
				// Verify the signature
				message := m.input
				signature := m.input2
				address := m.input3
				valid, err := VerifySignature(message, signature, address)
				if err != nil {
					m.content = fmt.Sprintf("Error verifying signature: %v", err)
					return m, nil
				}
				if valid {
					m.content = "Signature is valid."
				} else {
					m.content = "Signature is invalid."
				}
				m.state = "display"
				return m, nil
			}
		case tea.KeyRunes:
			if m.step == 0 {
				m.input += string(msg.Runes)
			} else if m.step == 1 {
				m.input2 += string(msg.Runes)
			} else if m.step == 2 {
				m.input3 += string(msg.Runes)
			}
		case tea.KeyBackspace, tea.KeyDelete:
			if m.step == 0 && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			} else if m.step == 1 && len(m.input2) > 0 {
				m.input2 = m.input2[:len(m.input2)-1]
			} else if m.step == 2 && len(m.input3) > 0 {
				m.input3 = m.input3[:len(m.input3)-1]
			}
		}
	}
	return m, nil
}

func (m model) viewVerify() string {
	s := titleStyle.Render("Verify Signature") + "\n\n"
	if m.step == 0 {
		s += "Enter the message that was signed or press Esc to cancel:\n"
		s += inputStyle.Render(m.input)
		if m.content != "" {
			s += "\n\n" + m.content + "\n\nPress Enter to continue..."
		}
	} else if m.step == 1 {
		s += "Enter the signature (in hex format) or press Esc to cancel:\n"
		s += inputStyle.Render(m.input2)
		if m.content != "" {
			s += "\n\n" + m.content + "\n\nPress Enter to continue..."
		}
	} else if m.step == 2 {
		s += "Enter the Ethereum address of the signer or press Esc to cancel:\n"
		s += inputStyle.Render(m.input3)
		if m.content != "" {
			s += "\n\n" + m.content + "\n\nPress Enter to continue..."
		}
	}
	return s
}

func SignMessage(privateKeyHex, message string) (string, error) {
	// Convert private key from hex string to ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Hash the message
	msgHash := sha256.Sum256([]byte(message))

	// Sign the hash
	signatureBytes, err := crypto.Sign(msgHash[:], privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	// Return the hex-encoded signature
	signatureHex := hexutil.Encode(signatureBytes)
	return signatureHex, nil
}

func VerifySignature(message, signatureHex, addressHex string) (bool, error) {
	// Validate the address
	if !common.IsHexAddress(addressHex) {
		return false, fmt.Errorf("invalid Ethereum address")
	}
	address := common.HexToAddress(addressHex)

	// Decode the signature
	signatureBytes, err := hexutil.Decode(signatureHex)
	if err != nil {
		return false, fmt.Errorf("invalid signature format")
	}
	if len(signatureBytes) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}

	// Hash the message
	msgHash := sha256.Sum256([]byte(message))

	// Remove recovery ID (last byte)
	sigPublicKey, err := crypto.Ecrecover(msgHash[:], signatureBytes)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Convert to ECDSA public key
	pubKey, err := crypto.UnmarshalPubkey(sigPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	// Generate address from public key
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Compare recovered address with provided address
	if recoveredAddr == address {
		return true, nil
	}
	return false, nil
}

func (m model) viewDisplay() string {
	s := m.content + "\n\nPress Enter to return to menu..."
	return s
}

func formatFarcasterData(fname string, result *airstack.FarcasterResponse) string {
	var sb strings.Builder

	social := result.Data.Socials.Social
	casts := result.Data.FarcasterCasts.Cast

	sb.WriteString(fmt.Sprintf("Results for Farcaster user '%s':\n\n", fname))

	if len(social) > 0 {
		s := social[0]
		sb.WriteString("Profile Information:\n")
		sb.WriteString(fmt.Sprintf("Profile Name   : %s\n", s.ProfileName))
		sb.WriteString(fmt.Sprintf("Follower Count : %d\n", s.FollowerCount))
		sb.WriteString(fmt.Sprintf("Following Count: %d\n", s.FollowingCount))
		sb.WriteString(fmt.Sprintf("FarScore       : %.2f\n", s.FarcasterScore.FarScore))
		sb.WriteString("\n")
	} else {
		sb.WriteString("No profile information found.\n\n")
	}

	if len(casts) > 0 {
		sb.WriteString("Recent Casts:\n")
		for i, cast := range casts {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, cast.Text))
		}
	} else {
		sb.WriteString("No recent casts found.\n")
	}

	return sb.String()
}
