package algosui

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/tui/createalgoui"
)

type model struct {
	user          string
	table         table.Model
	previousModel tea.Model
}

func AlgosNew(user string, previousModel tea.Model) model {
	algos, err := models.GetAlgos(user)
	if err != nil {
		log.Fatal(err)
	}

	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Ticket", Width: 10},
		{Title: "Status", Width: 10},
		{Title: "Performance", Width: 15},
	}

	rows := make([]table.Row, 0, 10)

	for _, algo := range algos {
		stats, err := models.GetStatsById2(algo.Id)
		if err != nil {
			log.Fatal(err)
		}
		stats_string := strconv.FormatFloat(stats.AvgReturnPerMonth*100, 'f', 4, 64) + " / mo"
		rows = append(rows, table.Row{algo.Name, algo.BaseAsset + algo.QuoteAsset, algo.State, stats_string})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(algos)+1),
	)

	return model{user: user, table: t, previousModel: previousModel}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m.previousModel, nil
		case tea.KeyCtrlN:
			return createalgoui.CreatealgoNew(m.user, m), nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case createalgoui.CreateAlgoMsg:
		if msg == createalgoui.UpdateAlgos {
			m.updateAlgosList()
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.table.View()
}
