package tts

import (
	"context"
	"testing"
	"time"

	"project_for_tmk_04_06/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----- 单元测试 -----
// TestSiliconFlowTTS_Interface 验证接口约束
func TestSiliconFlowTTS_Interface(t *testing.T) {
	var _ TTS = (*SiliconFlowTTS)(nil)
}

// TestSiliconFlowTTS_GetVoice 验证 GetVoiceForLang 返回正确的声音 ID。
// 原理：表驱动测试（Table-Driven Test）是 Go 的最佳实践。
// 用一个 slice of struct 描述输入和期望输出，循环执行，易于扩展。
func TestSiliconFlowTTS_GetVoice(t *testing.T) {
	ttsClient := NewSiliconFlowTTS()
	// 所有语言应该返回正确的 CosyVoice2 格式声音 ID
	voice := ttsClient.GetVoiceForLang("zh")
	assert.Equal(t, "FunAudioLLM/CosyVoice2-0.5B:alex", voice)
}

// ----- 集成测试 -----
// TestSiliconFlowTTS_Synthesize_Integration 调用 SiliconFlow TTS API。
// 验证返回的 MP3 数据非空，证明整条"文字→AI→音频字节"链路正常。
func TestSiliconFlowTTS_Synthesize_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test; use -run Integration to run")
	}

	_ = config.Load()

	ttsClient := NewSiliconFlowTTS()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 测试中文
	data, err := ttsClient.Synthesize(ctx, "你好，这是一段测试语音。", "zh")
	require.NoError(t, err, "Chinese TTS should not return error")
	assert.Greater(t, len(data), 1000, "Chinese audio data should have meaningful size")

	// 测试日文 —— 验证多语言支持
	dataJa, err := ttsClient.Synthesize(ctx, "こんにちは、テストです。", "ja")
	require.NoError(t, err, "Japanese TTS should not return error")
	assert.Greater(t, len(dataJa), 1000, "Japanese audio data should have meaningful size")
}
