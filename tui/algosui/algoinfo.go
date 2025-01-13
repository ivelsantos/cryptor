package algosui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivelsantos/cryptor/models"
)

type algoInfoModel struct {
	Algo models.Algor
	Text textarea.Model
}

func algoInfoNew(botid int) tea.Model {
	model := algoInfoModel{}
	model.Algo, _ = models.GetAlgoById(botid)

	model.Text = textarea.New()
	model.Text.ShowLineNumbers = true
	model.Text.SetValue(model.Algo.Buycode)
	model.Text.MaxHeight = 15
	model.Text.SetHeight(15)
	model.Text.SetWidth(50)

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
	style := lipgloss.NewStyle().Bold(true).AlignHorizontal(lipgloss.Center)
	s := lipgloss.JoinVertical(0.5, style.Render("Code\n"), a.Text.View())
	return s
}
