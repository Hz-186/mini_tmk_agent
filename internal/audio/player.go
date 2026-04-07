package audio

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/hajimehoshi/go-mp3"
)

// PlayMP3 plays the given MP3 byte slice using malgo sequentially (blocking).
func PlayMP3(mp3Data []byte) error {
	decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	if err != nil {
		return err
	}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {})
	if err != nil {
		return err
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = uint32(decoder.SampleRate())
	deviceConfig.Alsa.NoMMap = 1

	onSendFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		n, err := decoder.Read(pOutputSample)
		if err != nil && err != io.EOF {
			log.Printf("MP3 decode error: %v", err)
		}
		if n < len(pOutputSample) {
			for i := n; i < len(pOutputSample); i++ {
				pOutputSample[i] = 0
			}
		}
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSendFrames,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		return err
	}
	defer device.Uninit()

	if err := device.Start(); err != nil {
		return err
	}
	duration := float64(decoder.Length()) / float64(decoder.SampleRate()*4) // 4 bytes per frame (16-bit stereo)

	importTime := make(chan struct{})
	go func() {
		importTime <- struct{}{}
	}()

	<-importTime
	sleepTime := float64(time.Second) * duration
	time.Sleep(time.Duration(sleepTime))

	return nil
}
