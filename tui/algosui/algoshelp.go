package algosui

import (
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ivelsantos/cryptor/models"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Status key.Binding
	Delete key.Binding
	Create key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Status, k.Create, k.Delete},
		{k.Help, k.Back, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Status: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "change status"),
	),
	Create: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new algo"),
	),
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete algo"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "previous"),
	),
}

type refreshMsg struct{}

type sessionState uint

const (
	tableView sessionState = iota
	algoInfoView
)

func DoRefresh() tea.Msg {
	time.Sleep(5 * time.Second)
	return refreshMsg{}
}

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
		stats_string := strconv.FormatFloat(stats.AvgReturnPerDay*100, 'f', 3, 64) + "% / day"
		if stats_string == "0.000% / day" {
			stats_string = "..."
		}

		id_string := strconv.Itoa(algo.Id)
		rows = append(rows, table.Row{id_string, algo.Name, algo.BaseAsset + algo.QuoteAsset, algo.State, stats_string})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(algos)+2),
	)

	return t
}
