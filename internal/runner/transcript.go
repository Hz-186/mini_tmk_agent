package runner

import (
	"context"
	"fmt"
	"os"
	"project_for_tmk_04_06/internal/ai/asr"

	"github.com/pterm/pterm"
)

// 离线转写执行器
type TranscriptRunner struct {
	asrClient asr.ASR
}

func NewTranscriptRunner() *TranscriptRunner {
	return &TranscriptRunner{
		asrClient: asr.NewSiliconFlowASR(),
	}
}

func (r *TranscriptRunner) Run(ctx context.Context, audioFile string, outputFile string) error {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Transcribing audio file: %s...", audioFile))
	// 加载动画

	// Verify file exists
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		spinner.Fail("Audio file does not exist")
		return fmt.Errorf("audio file not found: %s", audioFile)
	}

	text, err := r.asrClient.Transcribe(ctx, audioFile)
	if err != nil {
		spinner.Fail("Transcription failed")
		return fmt.Errorf("API request failed: %w", err)
	}

	spinner.Success("Transcription finished successfully")
	err = os.WriteFile(outputFile, []byte(text), 0644) // 你可以读写这个文件，别人只能读这个文件
	if err != nil {
		pterm.Error.Printf("Failed to save transcript to %s\n", outputFile)
		return err
	}
	pterm.Success.Printf("Saved transcript to: %s\n", outputFile)
	return nil
}
