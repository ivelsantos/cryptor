package algosui

import (
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/ivelsantos/cryptor/models"
)

func (m *model) deleteAlgo(id int) {
	err := models.DeleteAlgo(id, m.user)
	if err != nil {
		panic(err)
	}
}

func (m *model) updateAlgosList() {
	m.table = getAlgosTable(m.user)
}

func getAlgosTable(user string) table.Model {
	algos, err := models.GetAlgos(user)
	if err != nil {
		log.Fatal(err)
	}

	columns := []table.Column{
		{Title: "Id", Width: 5},
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
		id_string := strconv.Itoa(algo.Id)
		rows = append(rows, table.Row{id_string, algo.Name, algo.BaseAsset + algo.QuoteAsset, algo.State, stats_string})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(algos)+1),
	)

	return t
}
