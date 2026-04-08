package audio

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestComputeRMS_Silence 验证静音信号返回接近 0 的 RMS 值。

func TestComputeRMS_Silence(t *testing.T) {
	// 构造全 0 的 PCM 数据 = 绝对静音
	silentPCM := make([]byte, 3200) // 0.1 秒 16kHz mono

	rms := RMS(silentPCM)
	assert.Equal(t, 0.0, rms, "Silent signal should have RMS of 0")
}

// TestComputeRMS_MaxAmplitude 验证最大振幅信号返回接近 1.0 的 RMS 值。
func TestComputeRMS_MaxAmplitude(t *testing.T) {
	buf := make([]byte, 400) // 100 个采样
	for i := 0; i < len(buf); i += 2 {
		buf[i] = 0xFF   // 低字节
		buf[i+1] = 0x7F // 高字节 → int16(0x7FFF) = 32767
	}
	rms := RMS(buf)
	assert.InDelta(t, 1.0, rms, 0.001, "Max amplitude signal should have RMS ≈ 1.0")
}

// TestComputeRMS_SineWave 验证正弦波信号返回 RMS ≈ Amplitude/√2。

// 对于振幅为 A 的纯正弦波 x(t) = A·sin(t)：
//
//	RMS = A / sqrt(2) ≈ 0.707 × A
//
// 这是因为 sin²(t) 的平均值是 0.5，所以 sqrt(0.5) = 1/sqrt(2)。
func TestComputeRMS_SineWave(t *testing.T) {
	numSamples := 1600 // 100ms 的 16kHz 音频
	buf := make([]byte, numSamples*2)

	amplitude := 16000.0 // 约为 32768 的一半
	for i := 0; i < numSamples; i++ {
		sample := int16(amplitude * math.Sin(2*math.Pi*float64(i)/float64(numSamples)))
		buf[i*2] = byte(sample)
		buf[i*2+1] = byte(sample >> 8)
	}

	rms := RMS(buf)
	expectedRMS := (amplitude / 32768.0) / math.Sqrt2
	assert.InDelta(t, expectedRMS, rms, 0.01, "Sine wave RMS should be ≈ A/√2")
}

// TestComputeRMS_Empty 验证空数据的边界安全。
func TestComputeRMS_Empty(t *testing.T) {
	rms := RMS([]byte{})
	assert.Equal(t, 0.0, rms, "Empty buffer should return 0 without panic")
}
