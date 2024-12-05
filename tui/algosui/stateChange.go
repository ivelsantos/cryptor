package algosui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
)

type stateChangeModel struct {
	algoId        int
	previousModel model
	rowPosition   int
}

func stateChangeNew(algoId int, rowPosition int, previousModel model) stateChangeModel {
	return stateChangeModel{previousModel: previousModel, algoId: algoId, rowPosition: rowPosition}
}

func (m stateChangeModel) Init() tea.Cmd {
	return nil
}

func (m stateChangeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.previousModel, nil
		case "W":
			models.UpdateAlgoState("waiting", m.algoId, m.previousModel.user)
			m.previousModel.updateAlgosList()
			m.previousModel.table.SetCursor(m.rowPosition)
			return m.previousModel, nil
		case "T":
			models.UpdateAlgoState("testing", m.algoId, m.previousModel.user)
			m.previousModel.updateAlgosList()
			m.previousModel.table.SetCursor(m.rowPosition)
			return m.previousModel, nil
		case "L":
			models.UpdateAlgoState("live", m.algoId, m.previousModel.user)
			m.previousModel.updateAlgosList()
			m.previousModel.table.SetCursor(m.rowPosition)
			return m.previousModel, nil
		}
	}
	return m, nil
}

func (m stateChangeModel) View() string {
	var b strings.Builder
	b.WriteString(m.previousModel.View())
	b.WriteString("\n\n\n\n")
	b.WriteString("States: (W) Waiting\t(T) Testing\t(L) Live\n")

	return b.String()
}
