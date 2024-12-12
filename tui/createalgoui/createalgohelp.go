package createalgoui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Previous key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Help, k.Previous, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Previous: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "previous"),
	),
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	textareaNewModel, cmd := m.textarea.Update(msg)
	m.textarea = textareaNewModel
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *model) verifySymbol() error {
	sym := m.inputs[1].Value() + m.inputs[2].Value()

	algos, err := models.GetAllAlgos()
	if err != nil {
		return err
	}
	for _, algo := range algos {
		if (algo.BaseAsset + algo.QuoteAsset) == sym {
			return nil
		}
	}

	account, err := models.GetAccountByName(m.user)
	if err != nil {
		return err
	}

	symbols, err := functions.GetSymbols(account.ApiKey, account.SecretKey)
	if err != nil {
		return err
	}

	for _, symbol := range symbols {
		if symbol.Symbol == sym {
			return nil
		}
	}

	return fmt.Errorf("Symbol does not exists\n")
}

func (m *model) verifyAlgo() error {
	algo := models.Algor{Created: time.Now().Unix(), State: "verification"}
	algo.Owner = m.user
	algo.Name = m.inputs[0].Value()
	algo.BaseAsset = m.inputs[1].Value()
	algo.QuoteAsset = m.inputs[2].Value()
	algo.Buycode = m.textarea.Value()

	optAlgo := lang.GlobalStore("Algo", algo)

	_, err := lang.Parse("", []byte(algo.Buycode), optAlgo)
	if err != nil {
		return err
	}
	return nil
}
