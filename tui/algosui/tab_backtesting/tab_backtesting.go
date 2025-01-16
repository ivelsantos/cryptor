package tabbacktesting

import tea "github.com/charmbracelet/bubbletea"

type backState uint

const (
	configView backState = iota
	pendingView
	resultView
)

type TabBacktesting struct {
	Config tea.Model
	State  backState
}

func New_TabBacktesting() tea.Model {
	m := TabBacktesting{}
	m.State = configView
	m.Config = new_BackConfig()

	return m
}

func (t TabBacktesting) Init() tea.Cmd {
	return nil
}

func (t TabBacktesting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t TabBacktesting) View() string {
	switch t.State {
	case configView:
		return t.Config.View()
	case pendingView:
		return "pending"
	default:
		return "result"
	}
}
