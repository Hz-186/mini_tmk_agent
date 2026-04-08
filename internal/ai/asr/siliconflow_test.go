package asr

import (
	"context"
	"testing"
	"time"

	"project_for_tmk_04_06/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----- 单元测试 (Unit Tests) -----
// 原理：只验证代码结构与接口约定，不发任何网络请求，瞬间完成。

// TestASRInterface 验证 SiliconFlowASR 是否满足 ASR 接口约定。
// 这是 Go 接口"鸭子类型"的核心验证方式：
//
//	var _ ASR = (*SiliconFlowASR)(nil)
func TestASRInterface(t *testing.T) {
	var _ ASR = (*SiliconFlowASR)(nil)
}

// TestNewSiliconFlowASR_Init 验证构造函数能否正常初始化，不崩溃。
// 原理：Go 中的"表驱动测试"——即使环境变量未设置，函数也应安全返回。
func TestNewSiliconFlowASR_Init(t *testing.T) {
	_ = config.Load() // 加载配置（可能是空的 key）
	asr := NewSiliconFlowASR()
	assert.NotNil(t, asr, "NewSiliconFlowASR() should return a valid instance")
	assert.Equal(t, "FunAudioLLM/SenseVoiceSmall", asr.model, "Default model should be SenseVoiceSmall")
}

// ----- 集成测试 (Integration Tests) -----
//   原理：testing.Short() 是 Go 官方约定的跳过机制。
//   运行 go test -short 时跳过，不消耗网络配额。
//   运行 go test -run Integration 时才真正发送 HTTP 请求到 SiliconFlow API。

// TestSiliconFlowASR_Transcribe_Integration 真实调用 API，验证端到端 WAV→文字链路。

func TestSiliconFlowASR_Transcribe_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test; use -run Integration to run")
	}

	_ = config.Load()
	client := NewSiliconFlowASR()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用项目根目录的真实录音文件（如果有的话）
	text, err := client.Transcribe(ctx, "../../../../test_tts.mp3")
	require.NoError(t, err, "Should transcribe without error")
	assert.NotEmpty(t, text, "Transcribed text should not be empty")
}
