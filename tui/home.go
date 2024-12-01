package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/table"
	"github.com/ivelsantos/cryptor/models"
)

func home(user string) string {
	algos, err := models.GetAlgos(user)
	if err != nil {
		log.Fatal(err)
	}

	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Ticket", Width: 10},
		{Title: "Status", Width: 15},
	}

	rows := make([]table.Row, 0, 10)

	for _, algo := range algos {
		rows = append(rows, table.Row{algo.Name, algo.BaseAsset + algo.QuoteAsset, algo.State})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(algos)),
	)

	return t.View()
}
