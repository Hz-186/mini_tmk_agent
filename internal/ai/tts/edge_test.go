package tts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEdgeTTS_Interface(t *testing.T) {
	var _ TTS = (*EdgeTTS)(nil)
}

func TestEdgeTTS_GetVoice(t *testing.T) {
	ttsInst := NewEdgeTTS()

	vZh := ttsInst.GetVoiceForLang("zh")
	assert.Equal(t, "zh-CN-XiaoxiaoNeural", vZh)

	vEs := ttsInst.GetVoiceForLang("es")
	assert.Equal(t, "es-ES-ElviraNeural", vEs)
}

// WA
func TestEdgeTTS_Synthesize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network/integration test")
	}

	ttsInst := NewEdgeTTS()
	data, err := ttsInst.Synthesize(context.Background(), "你好，世界", "zh")
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.GreaterOrEqual(t, len(data), 1000, "Synthesized audio should have some size")
}
