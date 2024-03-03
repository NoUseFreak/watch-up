package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nousefreak/watch-up/internal/app/ui"
	"github.com/nousefreak/watch-up/internal/app/watchup"
	"github.com/sirupsen/logrus"
)

func runApp(appSettings watchup.AppSettings) {
	logrus.SetLevel(logrus.DebugLevel)
	tmp := os.TempDir()
	logPath := fmt.Sprintf("%s/watch-up.log", tmp)
	logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(logfile)

	bus := watchup.ChanBus{
		ChangeResults: make(chan watchup.WatchResult),
		WatchResults:  make(chan watchup.WatchResult),
		CodeStats:     make(chan watchup.WatchCodeStats),
		RequestStats:  make(chan watchup.RequestStats),
		Shutdown:      make(chan os.Signal, 1),
	}
	signal.Notify(bus.Shutdown, syscall.SIGINT, syscall.SIGTERM)

	c := watchup.Collector{}
	w := watchup.NewWatcher(appSettings.URL, appSettings.LoopDuration)

	go startApp(bus, appSettings)
	go c.Start(bus)
	go w.Start(bus)

	<-bus.Shutdown

	c.Stop()
	w.Stop()
}

func startApp(bus watchup.ChanBus, appSettings watchup.AppSettings) {
	if _, err := tea.NewProgram(ui.NewAppModel(bus, appSettings), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
	bus.Shutdown <- syscall.SIGTERM
}
