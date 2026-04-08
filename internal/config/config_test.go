package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoad_Defaults 验证配置加载时，默认值被正确填充。

func TestLoad_Defaults(t *testing.T) {
	// 清空可能影响测试的环境变量，保证测试隔离性
	os.Unsetenv("TMK_AI_SILICON_FLOW_KEY")

	err := Load()
	require.NoError(t, err, "Load() should not return error")

	assert.Equal(t,
		"https://api.siliconflow.cn/v1",
		AppConfig.AI.BaseURL,
		"Default BaseURL should be SiliconFlow API",
	)
}

// TestLoad_FromEnv 验证通过操作系统环境变量注入 API Key 的机制是否有效。
func TestLoad_FromEnv(t *testing.T) {
	const testKey = "sk-test-unit-test-key"
	os.Setenv("TMK_AI_SILICON_FLOW_KEY", testKey)
	t.Cleanup(func() {
		os.Unsetenv("TMK_AI_SILICON_FLOW_KEY")
	})

	err := Load()
	require.NoError(t, err)

	assert.Equal(t, testKey, AppConfig.AI.Key,
		"API key from environment variable should be correctly loaded")
}
