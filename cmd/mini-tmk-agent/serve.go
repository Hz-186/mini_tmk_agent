package main

import (
	"log/slog"
	"os"
	"project_for_tmk_04_06/internal/runner"
	webserver "project_for_tmk_04_06/internal/web"

	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "开始主服务",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Starting Web UI server")

		go webserver.StartServer(servePort)
		ctx := cmd.Context()
		streamRn := runner.NewStreamRunner()

		// Since we want Web UI, we just hook up the event bus inside the stream runner.

		if err := streamRn.Run(ctx, sourceLang, targetLang, ttsEnabled); err != nil {
			slog.Error("Stream mode failed", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(&sourceLang, "source-lang", "zh", "Source language")
	serveCmd.Flags().StringVar(&targetLang, "target-lang", "en", "Target language")
	serveCmd.Flags().BoolVar(&ttsEnabled, "tts", false, "Enable TTS (Text-to-Speech) output")
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Web UI Port")
	rootCmd.AddCommand(serveCmd)
}
