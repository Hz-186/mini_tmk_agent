package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rtcCmd = &cobra.Command{
	Use:   "rtc",
	Short: "WebRTC peer-to-peer commands",
}
var rtcHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Host a new P2P room",
	Run: func(cmd *cobra.Command, args []string) {
		manager := webrtc.NewRTCManager()
		if err := manager.Host(cmd.Context(), sourceLang, targetLang, ttsEnabled); err != nil {
			pterm.Error.Println("RTC Host failed:", err)
			os.Exit(1)
		}
	},
}

var rtcJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a new P2P room",
	Run: func(cmd *cobra.Command, args []string) {
		manager := webrtc.NewRTCManager()
		if err := manager.Join(cmd.Context(), args[0]); err != nil {
			pterm.Error.Println("RTC Join failed:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rtcHostCmd.Flags().StringVar(&sourceLang, "source-lang", "zh", "Source language")
	rtcHostCmd.Flags().StringVar(&targetLang, "target-lang", "en", "Target language")
	rtcHostCmd.Flags().BoolVar(&ttsEnabled, "tts", false, "Enable TTS (Text-to-Speech) output")
	rtcCmd.AddCommand(rtcHostCmd)
	rtcCmd.AddCommand(rtcJoinCmd)
	rootCmd.AddCommand(rtcCmd)
}
