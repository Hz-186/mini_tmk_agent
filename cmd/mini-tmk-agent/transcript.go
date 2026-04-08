package main

import (
	"log/slog"
	"os"
	"project_for_tmk_04_06/internal/runner"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/sqweek/dialog"
)

var (
	transcriptFile   string
	transcriptOutput string
)

var transcriptCmd = &cobra.Command{
	Use:   "transcript",
	Short: "Start transcript mode to process an audio file",
	Run: func(cmd *cobra.Command, args []string) {
		if transcriptFile == "" {
			var err error
			pterm.Info.Println("Waiting for you to select an audio file in the popup dialog...")
			transcriptFile, err = dialog.File().Filter("Audio Files", "mp3", "wav", "pcm").Title("Select an audio file to transcribe").Load()
			if err != nil {
				pterm.Warning.Println("File selection canceled.")
				os.Exit(0)
			}
		}

		slog.Info("Starting transcript mode", "file", transcriptFile, "output", transcriptOutput)

		ctx := cmd.Context()
		rn := runner.NewTranscriptRunner()
		if err := rn.Run(ctx, transcriptFile, transcriptOutput); err != nil {
			slog.Error("Transcript execution failed", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	transcriptCmd.Flags().StringVar(&transcriptFile, "file", "", "Audio file to transcript (If empty, a file picker GUI will open)")
	transcriptCmd.Flags().StringVar(&transcriptOutput, "output", "transcript.txt", "Destination file path")
	rootCmd.AddCommand(transcriptCmd)
}
