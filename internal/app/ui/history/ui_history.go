package history

import (
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nousefreak/watch-up/internal/app/ui/style"
	"github.com/nousefreak/watch-up/internal/app/watchup"
)

func New(width, height int) HistoryViewModel {
	tableColumns := []table.Column{
		{Title: "Timestamp", Width: 27},
		{Title: "Code", Width: 10},
		{Title: "Duration", Width: 10},
	}

	t := table.New(
		table.WithColumns(tableColumns),
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

	return HistoryViewModel{
		Entries: []watchup.WatchResult{},
		Limit:   1000,
		Width:   width,
		Height:  height,
		table:   t,
		columns: tableColumns,
	}
}

type HistoryViewModel struct {
	table   table.Model
	columns []table.Column

	Entries []watchup.WatchResult
	Limit   int

	Width  int
	Height int
}

func (l *HistoryViewModel) View() string {
	data := []table.Row{}
	for i, entry := range l.Entries {
		duration := ""
		if i < len(l.Entries)-1 {
			duration = l.Entries[i+1].DeltaTime.Truncate(time.Second).String()
		} else {
			duration = time.Since(entry.Time).Truncate(time.Second).String()
		}
		data = append(data, table.Row{
			entry.Time.Format(time.RFC3339),
			watchup.FormatStatusCode(entry.StatusCode),
			duration,
		})
	}

	l.table.SetWidth(l.Width - 2)
	l.table.SetHeight(l.Height - 2)
	l.columns[2].Width = l.table.Width() - l.columns[0].Width - l.columns[1].Width - 6
	l.table.SetColumns(l.columns)

	l.table.SetRows(data)
	l.table.GotoBottom()

	return style.FocusedBorderStyle.Render(l.table.View())
}

func (l *HistoryViewModel) AddEntry(entry watchup.WatchResult) {
	l.Entries = append(l.Entries, entry)
	if len(l.Entries) > l.Limit {
		l.Entries = l.Entries[len(l.Entries)-l.Limit:]
	}
}

func (l HistoryViewModel) Update(msg tea.Msg) (HistoryViewModel, tea.Cmd) {
	l.table.SetWidth(l.Width)
	l.table.SetHeight(l.Height)
	// l.viewport.Width = l.Width
	// l.viewport.Height = l.Height

	var cmd tea.Cmd
	// l.viewport, cmd = l.viewport.Update(msg)
	l.table, cmd = l.table.Update(msg)
	return l, cmd
}
