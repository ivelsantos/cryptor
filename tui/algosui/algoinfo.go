package algosui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
)

type algoInfoModel struct {
	Algo models.Algor
}

func algoInfoNew(botid int) tea.Model {
	model := algoInfoModel{}
	model.Algo, _ = models.GetAlgoById(botid)

	return model
}

func (a algoInfoModel) Init() tea.Cmd {
	return nil
}

func (a algoInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return nil, cmd
}

func (a algoInfoModel) View() string {
	return a.Algo.Buycode
}
