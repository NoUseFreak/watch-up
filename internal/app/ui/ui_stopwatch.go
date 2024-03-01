package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
)

func newStopwatch() stopwatch.Model {
	s := stopwatch.NewWithInterval(time.Second)

	return s
}
