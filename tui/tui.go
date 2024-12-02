package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Tui() {

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type mainModel struct {
	models []tea.Model
	user   string
}

var main mainModel

func initialModel() mainModel {
	models := make([]tea.Model, 0, 5)
	models = append(models, loginNew())
	main.models = models
	return main
}

func popModel() tea.Model {
	if len(main.models) > 1 {
		main.models = main.models[:len(main.models)-1]
	}
	return main
}

func insertModel(model tea.Model) tea.Model {
	main.models = append(main.models, model)
	return main
}

func (m mainModel) Init() tea.Cmd {
	// return m.model.Init()
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return popModel(), nil
		}
	}

	return m.models[len(m.models)-1].Update(msg)
}

func (m mainModel) View() string {
	return m.models[len(m.models)-1].View()
}
