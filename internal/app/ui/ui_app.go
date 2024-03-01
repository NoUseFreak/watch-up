package ui

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nousefreak/watch-up/internal/app/watchup"
)

const (
	helpHeight = 5
)

var (
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	halfWidthStyle = lipgloss.NewStyle().
			Width(50)

	keyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})

	descStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})
)

type keymap = struct {
	pause, quit key.Binding
}

type model struct {
	width  int
	height int
	keymap keymap

	appSettings watchup.AppSettings

	bus watchup.ChanBus

	help       help.Model
	history    HistoryViewModel
	table      table.Model
	stopwatch  stopwatch.Model
	avgStat    StatBox
	jitterStat StatBox
	countStat  StatBox
}

func New(bus watchup.ChanBus, appSettings watchup.AppSettings) model {
	m := model{
		bus:         bus,
		appSettings: appSettings,
		history:     NewHistoryViewModel(20, 20),
		table:       newTable(),
		stopwatch:   newStopwatch(),
		avgStat:     NewStatBox("Average"),
		jitterStat:  NewStatBox("Jitter"),
		countStat:   NewStatBox("Count"),
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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.stopwatch.Start(),
		waitForCodeStats(m.bus),
		waitForChangeResults(m.bus),
		waitForRequestStats(m.bus),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		data := []table.Row{}
		keys := make([]int, 0, len(msg))
		for k := range msg {
			keys = append(keys, k)
		}

		sort.Ints(keys)

		for _, code := range keys {
			data = append(data, table.Row{
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

func (m *model) sizeInputs() {
	col := m.width / 12

	sBoxHeight := 5

	m.history.Width = col * 6
	m.history.Height = (m.height - helpHeight)
	m.table.SetWidth(col*6 - 2)
	tableColumns[1].Width = m.table.Width() - 6 - tableColumns[0].Width - tableColumns[2].Width
	m.table.SetColumns(tableColumns)
	m.table.SetHeight(m.height - helpHeight - sBoxHeight - 2)
	halfWidthStyle.Width(col * 6)

	m.avgStat.Width = col*2 - 2
	m.avgStat.Height = sBoxHeight - 2
	m.jitterStat.Width = col*2 - 2
	m.jitterStat.Height = sBoxHeight - 2
	m.countStat.Width = col*2 - 2
	m.countStat.Height = sBoxHeight - 2
}

func (m *model) updateKeybindings() {
	// m.keymap.add.SetEnabled(len(m.inputs) < maxInputs)
	// m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

func (m model) View() string {
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.pause,
		m.keymap.quit,
	})

	stopwatch := fmt.Sprintf(
		"%s %-10s",
		descStyle.Render("Elapsed:"),
		keyStyle.Render(m.stopwatch.View()),
	)

	titleString := halfWidthStyle.Align(lipgloss.Left).Render(
		keyStyle.Render(" [Watch Up] "+m.appSettings.URL),
	) + halfWidthStyle.Align(lipgloss.Right).Render(
		descStyle.Render(" Loop: ")+
			keyStyle.Render(m.appSettings.LoopDuration.String()),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		titleString,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.history.View(),
			lipgloss.JoinVertical(
				lipgloss.Top,
				focusedBorderStyle.Render(m.table.View()),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					focusedBorderStyle.Render(m.avgStat.View()),
					focusedBorderStyle.Render(m.jitterStat.View()),
					focusedBorderStyle.Render(m.countStat.View()),
				),
			),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			halfWidthStyle.Align(lipgloss.Left).Render(help),
			halfWidthStyle.Align(lipgloss.Right).Render(stopwatch),
		),
	)
}
