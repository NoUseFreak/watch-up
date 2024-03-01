package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nousefreak/watch-up/internal/app/watchup"
	"github.com/spf13/cobra"
)

var (
	appSettings watchup.AppSettings
)

var rootCmd = &cobra.Command{
	Use:   "watch-up URL",
	Short: "watch-up is a simple tool to monitor the uptime of a website",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appSettings.URL = args[0]
		runApp(appSettings)
	},
}

func main() {
	rootCmd.Flags().DurationVarP(&appSettings.LoopDuration, "loop", "l", 500*time.Millisecond, "Duration between each request")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
