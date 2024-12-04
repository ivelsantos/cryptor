package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/tui/algosui"
)

func Tui() {

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
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

	m := model{users: users}

	return m
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
			newModel := algosui.AlgosNew(m.users[m.cursor], m)
			return newModel, nil
		}
	}
	return m, nil
}

func (m model) View() string {
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
