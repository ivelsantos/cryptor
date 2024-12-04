package algosui

import (
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
)

type algosModel struct {
	user  string
	table table.Model
}

type algosMsg table.Model

func algosNew(user string) tea.Model {
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

	return algosModel{user: user, table: t}
}

func (m algosModel) Init() tea.Cmd {
	return m.updateAlgosList()
}

func (m algosModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return main.popModel(), nil
		case tea.KeyCtrlN:
			return main.insertModel(createalgoNew(m.user)), nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case algosMsg:
		m.table = table.Model(msg)
		return m, nil
	}

	return m, nil
}

func (m algosModel) View() string {
	return m.table.View()
}

func (m *algosModel) updateAlgosList() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(10 * time.Second)
		algos, err := models.GetAlgos(m.user)
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
		return algosMsg(t)
	}
}
