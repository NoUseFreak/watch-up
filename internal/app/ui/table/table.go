package table

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var Columns = []table.Column{
	{Title: "Code", Width: 10},
	{Title: "Name", Width: 30},
	{Title: "Duration", Width: 10},
}

func New() table.Model {
	t := table.New(
		table.WithColumns(Columns),
		table.WithFocused(false),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.NoColor{})
	t.SetStyles(s)

	return t
}
