package algosui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ivelsantos/cryptor/models"
)

type algoInfoModel struct {
	Tabs       []string
	TabContent []string
	activeTab  int
	Algo       models.Algor
	Text       textarea.Model
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

	model.Tabs = []string{"Code", "Live Performance", "Testing Performance", "Backtesting"}
	model.TabContent = []string{model.Text.View(), "Blush Tab", "Eye Shadow Tab", "Mascara Tab"}

	return model
}

func (a algoInfoModel) Init() tea.Cmd {
	return nil
}

func (a algoInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "right", "l", "n", "tab":
			a.activeTab = min(a.activeTab+1, len(a.Tabs)-1)
			return a, nil
		case "left", "h", "p", "shift+tab":
			a.activeTab = max(a.activeTab-1, 0)
			return a, nil
		}
	}

	return a, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (a algoInfoModel) View() string {
	// style := lipgloss.NewStyle().Bold(true).AlignHorizontal(lipgloss.Center)
	// s := lipgloss.JoinVertical(0.5, style.Render("Code\n"), a.Text.View())
	// return s
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range a.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(a.Tabs)-1, i == a.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(a.TabContent[a.activeTab]))
	return docStyle.Render(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
