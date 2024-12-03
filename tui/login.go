package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
	"log"
)

type loginModel struct {
	user   string
	users  []string
	cursor int
}

func loginNew() tea.Model {
	users := make([]string, 0, 5)
	accounts, err := models.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}

	for _, account := range accounts {
		users = append(users, account.Name)
	}

	m := loginModel{users: users}

	return m
}

func (m loginModel) Init() tea.Cmd {
	return nil
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
			return insertModel(algosNew(m.user)), nil
		}
	}
	return main, nil
}

func (m loginModel) View() string {
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
