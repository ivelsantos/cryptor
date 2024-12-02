package tui

import tea "github.com/charmbracelet/bubbletea"

type createalgoModel struct{}

func createalgoNew() tea.Model {
	return createalgoModel{}
}

func (m createalgoModel) Init() tea.Cmd {
	return nil
}

func (m createalgoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return main, nil
}

func (m createalgoModel) View() string {
	return "Creating a new Algo..."
}
