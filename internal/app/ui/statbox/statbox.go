package statbox

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatBox struct {
	Title string
	Value string

	Height int
	Width  int
}

func New(title string) StatBox {
	return StatBox{
		Title:  title,
		Value:  "0",
		Height: 5,
		Width:  20,
	}
}

func (s StatBox) Init() tea.Cmd {
	return nil
}

func (s StatBox) Update(msg tea.Msg) (StatBox, tea.Cmd) {
	return s, nil
}

func (s StatBox) View() string {
	title := lipgloss.NewStyle().Width(s.Width).Align(lipgloss.Center).Render(s.Title)
	value := lipgloss.NewStyle().Width(s.Width).Align(lipgloss.Center).Render(s.Value)

	return lipgloss.NewStyle().Height(s.Height).Align(lipgloss.Center).Render(lipgloss.JoinVertical(lipgloss.Top, title, "", value))
}
