package ui

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tablelib "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nousefreak/watch-up/internal/app/ui/history"
	"github.com/nousefreak/watch-up/internal/app/ui/statbox"
	"github.com/nousefreak/watch-up/internal/app/ui/style"
	"github.com/nousefreak/watch-up/internal/app/ui/table"
	"github.com/nousefreak/watch-up/internal/app/watchup"
)

const (
	helpHeight = 5
)

type keymap = struct {
	pause, quit key.Binding
}

type appModel struct {
	width  int
	height int
	keymap keymap

	appSettings watchup.AppSettings

	bus watchup.ChanBus

	help       help.Model
	history    history.HistoryViewModel
	table      tablelib.Model
	stopwatch  stopwatch.Model
	avgStat    statbox.StatBox
	jitterStat statbox.StatBox
	countStat  statbox.StatBox
}

// NewAppModel returns a new appModel.
func NewAppModel(bus watchup.ChanBus, appSettings watchup.AppSettings) appModel {
	m := appModel{
		bus:         bus,
		appSettings: appSettings,
		history:     history.New(20, 20),
		table:       table.New(),
		stopwatch:   stopwatch.New(),
		avgStat:     statbox.New("Average"),
		jitterStat:  statbox.New("Jitter"),
		countStat:   statbox.New("Count"),
		help:        help.New(),
		keymap: keymap{
			pause: key.NewBinding(
				key.WithKeys(tea.KeySpace.String()),
				key.WithHelp("space", "pause"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
	m.updateKeybindings()
	return m
}

func waitForCodeStats(bus watchup.ChanBus) tea.Cmd {
	return func() tea.Msg {
		return <-bus.CodeStats
	}
}

func waitForChangeResults(bus watchup.ChanBus) tea.Cmd {
	return func() tea.Msg {
		return <-bus.ChangeResults
	}
}

func waitForRequestStats(bus watchup.ChanBus) tea.Cmd {
	return func() tea.Msg {
		return <-bus.RequestStats
	}
}

// Init initializes the appModel.
func (m appModel) Init() tea.Cmd {
	return tea.Batch(
		m.stopwatch.Start(),
		waitForCodeStats(m.bus),
		waitForChangeResults(m.bus),
		waitForRequestStats(m.bus),
	)
}

// Update updates the appModel.
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.pause):
			return m, m.stopwatch.Toggle()
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case watchup.WatchCodeStats:
		data := []tablelib.Row{}
		keys := make([]int, 0, len(msg))
		for k := range msg {
			keys = append(keys, k)
		}

		sort.Ints(keys)

		for _, code := range keys {
			data = append(data, tablelib.Row{
				watchup.FormatStatusCode(code),
				watchup.CodeToText(code),
				msg[code].Truncate(time.Second).String(),
			})
		}

		m.table.SetRows(data)
		return m, waitForCodeStats(m.bus)
	case watchup.WatchResult:
		m.history.AddEntry(msg)

		return m, waitForChangeResults(m.bus)
	case watchup.RequestStats:
		m.avgStat.Value = msg.AvgTime.Truncate(time.Millisecond).String()
		m.jitterStat.Value = msg.Jitter.Truncate(time.Millisecond).String()
		m.countStat.Value = fmt.Sprintf("%d", msg.TotalRequests)
		return m, waitForRequestStats(m.bus)
	}

	m.updateKeybindings()
	m.sizeInputs()

	var cmd tea.Cmd
	m.history, cmd = m.history.Update(msg)
	cmds = append(cmds, cmd)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)

	m.avgStat, cmd = m.avgStat.Update(msg)
	cmds = append(cmds, cmd)

	m.jitterStat, cmd = m.jitterStat.Update(msg)
	cmds = append(cmds, cmd)

	m.countStat, cmd = m.countStat.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *appModel) sizeInputs() {
	col := m.width / 12

	sBoxHeight := 5

	m.history.Width = col * 6
	m.history.Height = (m.height - helpHeight)
	m.table.SetWidth(col*6 - 2)
	table.Columns[1].Width = m.table.Width() - 6 - table.Columns[0].Width - table.Columns[2].Width
	m.table.SetColumns(table.Columns)
	m.table.SetHeight(m.height - helpHeight - sBoxHeight - 2)
	style.HalfWidthStyle.Width(col * 6)

	m.avgStat.Width = col*2 - 2
	m.avgStat.Height = sBoxHeight - 2
	m.jitterStat.Width = col*2 - 2
	m.jitterStat.Height = sBoxHeight - 2
	m.countStat.Width = col*2 - 2
	m.countStat.Height = sBoxHeight - 2
}

func (m *appModel) updateKeybindings() {
	// m.keymap.add.SetEnabled(len(m.inputs) < maxInputs)
	// m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

// View returns the appModel's view.
func (m appModel) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.pause,
		m.keymap.quit,
	})

	stopwatch := fmt.Sprintf(
		"%s %-10s",
		style.DescStyle.Render("Elapsed:"),
		style.KeyStyle.Render(m.stopwatch.View()),
	)

	titleString := style.HalfWidthStyle.Align(lipgloss.Left).Render(
		style.KeyStyle.Render(" [Watch Up] "+m.appSettings.URL),
	) + style.HalfWidthStyle.Align(lipgloss.Right).Render(
		style.DescStyle.Render(" Loop: ")+
			style.KeyStyle.Render(m.appSettings.LoopDuration.String()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		titleString,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.history.View(),
			lipgloss.JoinVertical(
				lipgloss.Top,
				style.FocusedBorderStyle.Render(m.table.View()),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					style.FocusedBorderStyle.Render(m.avgStat.View()),
					style.FocusedBorderStyle.Render(m.jitterStat.View()),
					style.FocusedBorderStyle.Render(m.countStat.View()),
				),
			),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			style.HalfWidthStyle.Align(lipgloss.Left).Render(help),
			style.HalfWidthStyle.Align(lipgloss.Right).Render(stopwatch),
		),
	)
}
