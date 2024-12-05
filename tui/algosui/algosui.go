package algosui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/tui/createalgoui"
)

type model struct {
	user          string
	table         table.Model
	previousModel tea.Model
}

func AlgosNew(user string, previousModel tea.Model) model {
	t := getAlgosTable(user)
	return model{user: user, table: t, previousModel: previousModel}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyEsc:
			return m.previousModel, nil

		case tea.KeyCtrlN:
			return createalgoui.CreatealgoNew(m.user, m), nil

		case tea.KeyCtrlS:
			id, _ := strconv.Atoi(m.table.SelectedRow()[0])
			return stateChangeNew(id, m.table.Cursor(), m), nil

		case tea.KeyCtrlD:
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

			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit
		}

	case createalgoui.CreateAlgoMsg:
		if msg == createalgoui.UpdateAlgos {
			m.updateAlgosList()
			return m, nil
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.table.View()
}
