package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
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

type createalgoModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	textarea   textarea.Model
	user       string
}

func createalgoNew(user string) tea.Model {
	m := createalgoModel{
		inputs: make([]textinput.Model, 3),
		user:   user,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Prompt = "Algo Name: "
			t.Placeholder = "..."
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "Base Asset: "
			t.Placeholder = "..."
			t.CharLimit = 5
		case 2:
			t.Prompt = "Quote Asset: "
			t.Placeholder = "..."
			t.CharLimit = 5
		}

		m.inputs[i] = t
	}

	ta := textarea.New()
	ta.Placeholder = "Code here..."
	ta.CharLimit = 500
	m.textarea = ta

	return m
}

func (m createalgoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m createalgoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			newMain := main.popModel()
			return newMain, nil

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

			if m.focusIndex == len(m.inputs) {
				cmds = append(cmds, m.textarea.Focus())
			} else {
				m.textarea.Blur()
			}

			return m, tea.Batch(cmds...)
		case "enter":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs)+1 {

				algo := models.Algor{Created: time.Now().Unix(), State: "waiting"}
				algo.Owner = m.user
				algo.Name = m.inputs[0].Value()
				algo.BaseAsset = m.inputs[1].Value()
				algo.QuoteAsset = m.inputs[2].Value()
				algo.Buycode = m.textarea.Value()

				err := models.InsertAlgo(algo)
				if err != nil {
					panic(err)
				}

				newMain := main.popModel()
				return newMain, nil
			}
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *createalgoModel) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m createalgoModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	b.WriteString("\n\n" + m.textarea.View() + "\n")

	button := &blurredButton
	if m.focusIndex == len(m.inputs)+1 {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
