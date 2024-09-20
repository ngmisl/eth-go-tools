// app.go

package main

import (
	"fmt"
	"os"
	"strings"

	"example.com/ethgotools/airstack"
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
	state    string
	// err      error
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
		return m.updateGenerate(msg)
	case "farcaster":
		return m.updateFarcaster(msg)
	case "display":
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			m.state = "menu"
			m.content = ""
			return m, nil
		}
		return m, nil
	default:
		return m, nil
	}
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

func (m model) updateGenerate(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					return fmt.Sprintf("Error: AIRSTACK_API_KEY not set.")
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
