package log

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nousefreak/watch-up/internal/app/ui/style"
)

func NewLogViewModel(width, height int) LogViewModel {
	vp := viewport.New(width, height)
	vp.Style = style.FocusedBorderStyle
	return LogViewModel{
		Entries:  []string{},
		Limit:    height,
		Width:    width,
		Height:   height,
		viewport: vp,
	}
}

type LogViewModel struct {
	viewport viewport.Model

	Entries []string
	Limit   int

	Width  int
	Height int
}

func (l *LogViewModel) View() string {
	l.viewport.SetContent(strings.Join(l.Entries, "\n"))

	return l.viewport.View()
}

func (l *LogViewModel) AddEntry(entry string) {
	l.Entries = append(l.Entries, strings.TrimSpace(entry))
	if len(l.Entries) > l.Limit {
		l.Entries = l.Entries[len(l.Entries)-l.Limit:]
	}
}

func (l LogViewModel) Update(msg tea.Msg) (LogViewModel, tea.Cmd) {
	l.viewport.Width = l.Width
	l.viewport.Height = l.Height

	var cmd tea.Cmd
	l.viewport, cmd = l.viewport.Update(msg)
	return l, cmd
}
