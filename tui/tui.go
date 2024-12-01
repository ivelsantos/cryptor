package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/ivelsantos/cryptor/models"

	tea "github.com/charmbracelet/bubbletea"
)

func Tui() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	user   string
	users  []string
	cursor int
}

func initialModel() model {
	users := make([]string, 0, 5)
	accounts, err := models.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range accounts {
		users = append(users, account.Name)
	}

	return model{
		users: users,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.users)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.user = m.users[m.cursor]
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.user != "" {
		page, err := home(m.user)
		if err != nil {
			log.Fatal(err)
		}
		return page
	}

	s := "\nChoose the user:\n\n"

	for i, choice := range m.users {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress ctrl+c to quit.\n"

	return s
}
