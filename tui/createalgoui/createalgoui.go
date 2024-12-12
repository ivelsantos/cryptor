package createalgoui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
)

type CreateAlgoMsg int

const (
	UpdateAlgos CreateAlgoMsg = iota
)

type model struct {
	focusIndex    int
	inputs        []textinput.Model
	textarea      textarea.Model
	user          string
	previousModel tea.Model
	keys          keyMap
	help          help.Model
	errAlgo       bool
	errSymbol     bool
}

func CreatealgoNew(user string, previousModel tea.Model) tea.Model {
	m := model{
		inputs:        make([]textinput.Model, 3),
		user:          user,
		previousModel: previousModel,
		keys:          keys,
		help:          help.New(),
		errAlgo:       false,
		errSymbol:     false,
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

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.errSymbol = false
	m.errAlgo = false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return m.previousModel, nil

			// Set focus to next input
		case "ctrl+h":
			m.help.ShowAll = !m.help.ShowAll
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

				errCode := m.verifySymbol()
				if errCode != nil {
					m.errSymbol = true
					return m, nil
				}

				errCode = m.verifyAlgo()
				if errCode != nil {
					m.errAlgo = true
					return m, nil
				}

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

				prevModel, cmd := m.previousModel.Update(UpdateAlgos)
				return prevModel, cmd
			}
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m model) View() string {
	// Inputs field
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

	var errCode string
	if m.errSymbol {
		errCode = "Symbol does not exists!\n"
	} else if m.errAlgo {
		errCode = "Code contain errors!\n"
	}

	s := b.String()

	helpView := m.help.View(m.keys)
	return s + "\n\n\n" + errCode + helpView
}
