package audio

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gen2brain/malgo"
)

type PcmChunk []byte

type SimpleRecorder struct {
	ctx                *malgo.AllocatedContext
	device             *malgo.Device // 代表麦克风设备本身
	sampleRate         uint32
	channels           uint16
	vadThreshold       float64
	silenceMaxDuration time.Duration
	bufferMutex        sync.Mutex
	pcmData            []byte
	lastVoiceTime      time.Time
	phraseChan         chan PcmChunk
	// 一个 channel，用来把完整的一句话发送出去
}

// 采样率  /  声道数
func NewSimpleRecorder(sampleRate uint32, channels uint16) (*SimpleRecorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {})
	if err != nil {
		return nil, err
	}

	return &SimpleRecorder{
		ctx:                ctx,
		sampleRate:         sampleRate,
		channels:           channels,
		phraseChan:         make(chan PcmChunk, 10),
		vadThreshold:       0.001, // 最低分贝
		silenceMaxDuration: 1 * time.Second,
	}, nil
}

func RMS(buffer []byte) float64 {
	if len(buffer) == 0 {
		return 0
	}
	var sumSquares float64
	// 16-bit sample = 2 bytes. We assume LittleEndian.
	totalSamples := len(buffer) / 2
	for i := 0; i < totalSamples; i++ {
		sample := float64(int16(buffer[i*2]) | int16(buffer[i*2+1])<<8)
		// normalize to -1.0 to 1.0 range
		normalized := sample / 32768.0
		sumSquares += normalized * normalized
	}
	return math.Sqrt(sumSquares / float64(totalSamples))
}

func (r *SimpleRecorder) Start(ctx context.Context) (<-chan PcmChunk, error) {
	//
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture) // malgo.Capture -> 录音
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = uint32(r.channels)
	deviceConfig.SampleRate = r.sampleRate
	// 录音设备配置

	r.lastVoiceTime = time.Now()

	// 写入录音的逻辑
	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		energy := RMS(pInputSamples) // 最大上限

		r.bufferMutex.Lock()
		defer r.bufferMutex.Unlock()

		if energy > r.vadThreshold {
			r.lastVoiceTime = time.Now()
			r.pcmData = append(r.pcmData, pInputSamples...)
		} else {
			if len(r.pcmData) > 0 {
				r.pcmData = append(r.pcmData, pInputSamples...)
			}

			minBytes := int(float32(r.sampleRate*2*uint32(r.channels)) * 0.3)
			if time.Since(r.lastVoiceTime) > r.silenceMaxDuration && len(r.pcmData) > minBytes {
				// Flush audio
				chunk := make([]byte, len(r.pcmData))

				copy(chunk, r.pcmData)

				select {
				case r.phraseChan <- chunk:
				default:
				}
				r.pcmData = make([]byte, 0)
				r.lastVoiceTime = time.Now()
			} else if time.Since(r.lastVoiceTime) > r.silenceMaxDuration {
				// 用户很久没说话，但刚刚录下来的声音不到 0.3 秒
				// 噪音。
				r.pcmData = r.pcmData[:0]
			}
		}
	}
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	// 初始化设备
	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		return nil, fmt.Errorf("InitDevice error: %w (请检查拔插麦克风硬件，或由于 Windows 隐私设置阻止了命令行访问麦克风)", err)
	}
	r.device = device

	// 开始录音
	if err := r.device.Start(); err != nil {
		return nil, fmt.Errorf("Start Device error: %w (请确保您的麦克风没有被其他程序独占)", err)
	}
	go func() {
		// 结束了 关闭channel
		<-ctx.Done()
		r.Stop()
	}()
	return r.phraseChan, nil
}

func (r *SimpleRecorder) Stop() {
	if r.device != nil {
		r.device.Uninit()
	}
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
	close(r.phraseChan)
}
