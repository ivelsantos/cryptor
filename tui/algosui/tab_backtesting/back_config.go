package tabbacktesting

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type BackConfig struct {
	Input textinput.Model
}

func new_BackConfig() tea.Model {
	m := BackConfig{}

	m.Input = textinput.New()
	m.Input.Prompt = "Days to backtest:"
	m.Input.CharLimit = 32

	return m
}

func (b BackConfig) Init() tea.Cmd {
	return nil
}

func (b BackConfig) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b BackConfig) View() string {
	return b.Input.View()
}
