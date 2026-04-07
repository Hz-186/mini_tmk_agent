package runner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"project_for_tmk_04_06/internal/ai/asr"
	"project_for_tmk_04_06/internal/ai/llm"
	"project_for_tmk_04_06/internal/ai/tts"
	"project_for_tmk_04_06/internal/audio"
	webserver "project_for_tmk_04_06/internal/web"
	"project_for_tmk_04_06/pkg"
	"strings"

	"github.com/pterm/pterm"
)

type StreamRunner struct {
	asrClient asr.ASR
	llmClient llm.LLM
	ttsClient tts.TTS
	Callback  func(source string, translation string)
}

func NewStreamRunner() *StreamRunner {
	return &StreamRunner{
		asrClient: asr.NewSiliconFlowASR(),
		llmClient: llm.NewSiliconFlowLLM(),
		ttsClient: tts.NewEdgeTTS(),
	}
}

func (r *StreamRunner) Run(ctx context.Context, sourceLang, targetLang string, enableTTS bool) error {
	pterm.DefaultHeader.WithFullWidth().Println("🎙️ Mini TMK Agent - Simultaneous Interpretation Stream")
	pterm.Info.Printf("Direction: %s -> %s | TTS Active: %v\n\n", sourceLang, targetLang, enableTTS)

	// 录音
	sampleRate := uint32(16000)
	channels := uint16(1)
	recorder, err := audio.NewSimpleRecorder(sampleRate, channels)
	if err != nil {
		return fmt.Errorf("failed to init recorder: %w", err)
	}

	// 开始录音
	phraseChan, err := recorder.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}
	pterm.Success.Println("Listening... Speak into the microphone OR type text directly here and press ENTER. (Press Ctrl+C to stop)")

	// 文字输入检测
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			if text != "" {
				pterm.Info.Printf("⌨️ [Manual Input Simulated]: %s\n", text)
				go r.processTextDirectly(ctx, text, sourceLang, targetLang, enableTTS)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case pcmChunk := <-phraseChan:
			if len(pcmChunk) == 0 {
				continue
			}

			// Generate valid WAV bytes in memory
			wavBytes, err := pkg.GenerateWAVBytes(sampleRate, channels, pcmChunk)
			// 把原始 PCM 数据加上 WAV 文件头，变成标准 WAV 格式
			// 因为AI的ASR接口只接受标准WAV格式，不接受裸PCM数据
			if err != nil {
				pterm.Error.Println("Audio format error:", err)
				continue
			}

			go r.processPhrase(ctx, wavBytes, sourceLang, targetLang, enableTTS)
		}
	}
	return nil
}

func (r *StreamRunner) processPhrase(ctx context.Context, audioWavData []byte, source string, target string, enableTTS bool) {
	spinner, _ := pterm.DefaultSpinner.Start("Analyzing voice...")
	text, err := r.asrClient.TranscribeBytes(ctx, audioWavData, source)
	if err != nil {
		spinner.Fail("ASR error: ", err)
		return
	}
	if text == "" {
		spinner.Warning("No speech detected clearly.")
		return
	}
	spinner.Success(fmt.Sprintf("[%s]: %s", source, text))
	// 完成语音的输入与转化文本

	streamPanel, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("[%s]: Translating...", target))
	translateChan, err := r.llmClient.TranslateStream(ctx, text, source, target)
	if err != nil {
		streamPanel.Fail("Translation failed: ", err)
		return
	}
	var fullTranslation string // 按照流读入，并且实时更新
	for token := range translateChan {
		fullTranslation += token
		streamPanel.UpdateText(fmt.Sprintf("[%s]: %s", target, fullTranslation))
	}
	streamPanel.Success(fmt.Sprintf("[%s]: %s", target, fullTranslation))
	// 完成翻译过程

	select { // select + default 非阻塞发送
	case webserver.EventBus <- webserver.TranslationEvent{
		Source:      text,
		Translation: fullTranslation,
	}:
	default: // don't block
	}

	if r.Callback != nil {
		r.Callback(text, fullTranslation)
	}

	if enableTTS && fullTranslation != "" {
		audioData, err := r.ttsClient.Synthesize(ctx, fullTranslation, target)
		if err != nil {
			pterm.Error.Println("TTS Synthesis failed:", err)
		} else if len(audioData) > 0 {
			pterm.Info.Println("🎵 [Playing TTS...]")
			if err := audio.PlayMP3(audioData); err != nil {
				pterm.Error.Println("Failed to play TTS:", err)
			}
		}
	}
}

func (r *StreamRunner) processTextDirectly(ctx context.Context, text string, source string, target string, enableTTS bool) {
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("[%s]: %s", source, text))
	spinner.Success()

	streamPanel, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("[%s]: Translating...", target))

	translateChan, err := r.llmClient.TranslateStream(ctx, text, source, target)
	if err != nil {
		streamPanel.Fail("Translation failed: ", err)
		return
	}

	var fullTranslation string
	for token := range translateChan {
		fullTranslation += token
		streamPanel.UpdateText(fmt.Sprintf("[%s]: %s", target, fullTranslation))
	}
	streamPanel.Success(fmt.Sprintf("[%s]: %s", target, fullTranslation))

	// Push to UI via EventBus if running in serve mode
	select {
	case webserver.EventBus <- webserver.TranslationEvent{
		Source:      text,
		Translation: fullTranslation,
	}:
	default: // don't block
	}

	if r.Callback != nil {
		r.Callback(text, fullTranslation)
	}

	if enableTTS && fullTranslation != "" {
		audioData, err := r.ttsClient.Synthesize(ctx, fullTranslation, target)
		if err == nil && len(audioData) > 0 {
			pterm.Info.Println("🎵 [Playing TTS...]")
			if err := audio.PlayMP3(audioData); err != nil {
				pterm.Error.Println("Failed to play TTS:", err)
			}
		}
	}
}
