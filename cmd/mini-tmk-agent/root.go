package main

import (
	"log/slog"
	logger "project_for_tmk_04_06/pkg"

	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "mini-tmk-agent",
	Short: "Mini TMK Agent is a CLI tool for simultaneous translation",
	Long:  `A powerful, purely Go-based CLI agent for real-time and file-based simultaneous interpretation.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			logger.EnableDebug()
		}
		slog.Debug("Debug logging enabled")
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
}
func Execute() error {
	return rootCmd.Execute()
}
