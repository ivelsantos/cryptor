package algosui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivelsantos/cryptor/tui/createalgoui"
	"golang.org/x/term"
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
	state         sessionState
}

func AlgosNew(user string, previousModel tea.Model) model {
	t := getAlgosTable(user)
	id, _ := strconv.Atoi(t.SelectedRow()[0])
	algoInfo := algoInfoNew(id)

	width, height, _ := term.GetSize(0)

	return model{
		user:          user,
		table:         t,
		algoInfo:      algoInfo,
		previousModel: previousModel,
		keys:          keys,
		help:          help.New(),
		width:         width,
		height:        height,
		state:         tableView,
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
		case "enter":
			if m.state == tableView {
				m.state = algoInfoView
			}
			return m, nil
		case "esc":
			if m.state == algoInfoView {
				m.state = tableView
				return m, nil
			}
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
	if m.state == tableView {
		m.table, cmd = m.table.Update(msg)
		// Update algoInfo view
		id, _ := strconv.Atoi(m.table.SelectedRow()[0])
		m.algoInfo = algoInfoNew(id)
	} else {
		m.algoInfo, cmd = m.algoInfo.Update(msg)
	}

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
		Width(int(float64(m.width) * 0.50)).
		Height(int(float64(m.height) * 0.9))

	var algoStyle = lipgloss.NewStyle().
		Width(int(float64(m.width) * 0.45)).
		MarginLeft(5)

	var helpStyle = lipgloss.NewStyle().
		Align(lipgloss.Bottom)

	var focusedStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("14"))

	var s string
	if m.state == tableView {
		s = lipgloss.JoinHorizontal(lipgloss.Top, focusedStyle.Render(tableStyle.Render(m.table.View())), algoStyle.Render(m.algoInfo.View()))
	} else {
		s = lipgloss.JoinHorizontal(lipgloss.Top, tableStyle.Render(m.table.View()), focusedStyle.Render(algoStyle.Render(m.algoInfo.View())))
	}
	s += "\n"
	s += helpStyle.Render(m.help.View(m.keys))

	return style.Render(s)
}
