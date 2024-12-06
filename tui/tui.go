package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/tui/algosui"
	"github.com/ivelsantos/cryptor/tui/createuser"
)

func Tui() {

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	users  []string
	cursor int
	keys   keyMap
	help   help.Model
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

	m := model{
		users: users,
		keys:  keys,
		help:  help.New(),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
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
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "enter":
			newModel := algosui.AlgosNew(m.users[m.cursor], m)
			return newModel, nil
		case "ctrl+n":
			newModel := createuser.CreateuserNew(m)
			return newModel, nil
		case "ctrl+d":
			n_users := len(m.users)
			m.deleteUser(m.users[m.cursor])

			m.updateUsers()

			if m.cursor == (n_users - 1) {
				m.cursor--
			}
			return m, nil
		}
	case createuser.CreateUserMsg:
		if msg == createuser.UpdateUsers {
			m.updateUsers()
			return m, nil
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

	helpView := m.help.View(m.keys)
	// height := 9 - strings.Count(s, "\n") - strings.Count(helpView, "\n")
	// if height < 0 {
	// 	height = 0
	// }

	// return s + strings.Repeat("\n", height) + helpView
	return s + "\n\n\n" + helpView
}
