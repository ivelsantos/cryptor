package algosui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivelsantos/cryptor/tui/createalgoui"
)

type model struct {
	user          string
	table         table.Model
	algoInfo      tea.Model
	previousModel tea.Model
	keys          keyMap
	help          help.Model
	height        int
	width         int
}

func AlgosNew(user string, previousModel tea.Model) model {
	t := getAlgosTable(user)
	id, _ := strconv.Atoi(t.SelectedRow()[0])
	algoInfo := algoInfoNew(id)
	return model{
		user:          user,
		table:         t,
		algoInfo:      algoInfo,
		previousModel: previousModel,
		keys:          keys,
		help:          help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case refreshMsg:
		index_row := m.table.Cursor()

		m.updateAlgosList()

		m.table.SetCursor(index_row)

		return m, DoRefresh
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.previousModel, nil

		case "?":
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case "ctrl+n":
			return createalgoui.CreatealgoNew(m.user, m), nil

		case "ctrl+s":
			id, _ := strconv.Atoi(m.table.SelectedRow()[0])
			return stateChangeNew(id, m.table.Cursor(), m), nil

		case "ctrl+d":
			id, _ := strconv.Atoi(m.table.SelectedRow()[0])
			m.deleteAlgo(id)
			n_rows := len(m.table.Rows())
			index_row := m.table.Cursor()

			m.updateAlgosList()

			if index_row != (n_rows - 1) {
				m.table.SetCursor(index_row)
			} else {
				m.table.GotoBottom()
			}

			// Update algoInfo view
			id, _ = strconv.Atoi(m.table.SelectedRow()[0])
			m.algoInfo = algoInfoNew(id)

			return m, nil

		case "ctrl+c":
			return m, tea.Quit
		}

	case createalgoui.CreateAlgoMsg:
		if msg == createalgoui.UpdateAlgos {
			index_row := m.table.Cursor()
			m.updateAlgosList()
			m.table.SetCursor(index_row)
			return m, DoRefresh
		}
	}
	m.table, cmd = m.table.Update(msg)

	// Update algoInfo view
	id, _ := strconv.Atoi(m.table.SelectedRow()[0])
	m.algoInfo = algoInfoNew(id)

	return m, cmd
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	var tableStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		BorderRight(true).
		Width(77).
		Height(20)

	var algoStyle = lipgloss.NewStyle().
		Width(66).
		MarginLeft(5)

	var helpStyle = lipgloss.NewStyle().
		Align(lipgloss.Bottom)

	s := lipgloss.JoinHorizontal(lipgloss.Top, tableStyle.Render(m.table.View()), algoStyle.Render(m.algoInfo.View()))
	s += "\n"
	s += helpStyle.Render(m.help.View(m.keys))

	return style.Render(s)
}
