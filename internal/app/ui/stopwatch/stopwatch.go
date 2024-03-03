package stopwatch

import (
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
)

func New() stopwatch.Model {
	s := stopwatch.NewWithInterval(time.Second)

	return s
}
