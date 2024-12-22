package createuser

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivelsantos/cryptor/models"
)

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

type CreateUserMsg int

const (
	UpdateUsers CreateUserMsg = iota
)

type model struct {
	focusIndex    int
	inputs        []textinput.Model
	previousModel tea.Model
	keys          keyMap
	help          help.Model
}

func CreateuserNew(previousModel tea.Model) tea.Model {
	m := model{
		inputs:        make([]textinput.Model, 5),
		previousModel: previousModel,
		keys:          keys,
		help:          help.New(),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Prompt = "Username: "
			t.Placeholder = "..."
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "Api key: "
			t.Placeholder = "..."
			t.CharLimit = 100
		case 2:
			t.Prompt = "Secret key: "
			t.Placeholder = "..."
			t.CharLimit = 100
		case 3:
			t.Prompt = "Testing api key: "
			t.Placeholder = "..."
			t.CharLimit = 100
		case 4:
			t.Prompt = "Testing secret key: "
			t.Placeholder = "..."
			t.CharLimit = 100
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return m.previousModel, nil

		case "ctrl+h":
			m.help.ShowAll = !m.help.ShowAll

			// Set focus to next input
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)+1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) + 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		case "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {

				user := models.Account{}
				user.Name = m.inputs[0].Value()
				user.ApiKey = m.inputs[1].Value()
				user.SecretKey = m.inputs[2].Value()

				err := models.InsertAccount(user)
				if err != nil {
					panic(err)
				}

				prevModel, cmd := m.previousModel.Update(UpdateUsers)
				return prevModel, cmd
				// return m.previousModel, nil
			}
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	s := b.String()
	helpView := m.help.View(m.keys)
	return s + "\n\n\n" + helpView
}
