package llm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSiliconFlowLLM_PromptBuilder(t *testing.T) {
	s := NewSiliconFlowLLM()
	
	messages := s.buildPrompt("hello", "en", "zh")
	require.Len(t, messages, 2)
	assert.Equal(t, "user", messages[1].Role)
	assert.Equal(t, "hello", messages[1].Content)
	assert.Contains(t, messages[0].Content, "from en to zh")
}

// We usually use httptest to mock the OpenAI client.
// Here we just test the instantiation and standard interfaces.
func TestLLMInterface(t *testing.T) {
	var _ LLM = (*SiliconFlowLLM)(nil)
}

// We want to skip actual network calls during standard test runs.
func TestSiliconFlowLLM_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	s := NewSiliconFlowLLM()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	res, err := s.Translate(ctx, "hello world", "en", "zh")
	require.NoError(t, err)
	assert.NotEmpty(t, res)
}
