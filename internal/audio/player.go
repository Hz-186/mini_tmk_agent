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

	// Playback callback
	onSendFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		n, err := decoder.Read(pOutputSample)
		if err != nil && err != io.EOF {
			log.Printf("MP3 decode error: %v", err)
		}
		// Fill with zero ifEOF
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

	// Calculate length to block until done
	duration := float64(decoder.Length()) / float64(decoder.SampleRate()*4) // 4 bytes per frame (16-bit stereo)

	importTime := make(chan struct{})
	go func() {
		// Blocking based on calculated duration + small padding since go-mp3 has an Exact length.
		// A cleaner approach checks decoded byte counts, but calculating duration is sufficient here.
		importTime <- struct{}{}
	}()
	<-importTime
	// Wait for the duration. Using timer to block till end of track.
	// Actually better is a small loop to check position, but we just sleep.
	sleepTime := float64(time.Second) * duration
	time.Sleep(time.Duration(sleepTime))

	return nil
}
