package main

import (
	"log/slog"
	"os"
	"project_for_tmk_04_06/internal/runner"

	"github.com/spf13/cobra"
)

var (
	sourceLang string
	targetLang string
	ttsEnabled bool
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: " ",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("Starting stream mode", "source", sourceLang, "target", targetLang, "tts", ttsEnabled)
		ctx := cmd.Context()
		streamRn := runner.NewStreamRunner()
		if err := streamRn.Run(ctx, sourceLang, targetLang, ttsEnabled); err != nil {
			slog.Error("Stream mode failed", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	streamCmd.Flags().StringVar(&sourceLang, "source-lang", "zh", "Source language (e.g. zh, en, es, ja)")
	streamCmd.Flags().StringVar(&targetLang, "target-lang", "en", "Target language (e.g. en, zh, es, ja)")
	streamCmd.Flags().BoolVar(&ttsEnabled, "tts", false, "Enable TTS (Text-to-Speech) output")
	rootCmd.AddCommand(streamCmd)
}
