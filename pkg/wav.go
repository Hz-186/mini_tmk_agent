package pkg

import (
	"bytes"
	"encoding/binary"
	"io"
)

// 把 pcm 转成 wav 文件，传给 api 接口使用
func GenerateWAVBytes(sampleRate uint32, channels uint16, pcmData []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	dataSize := uint32(len(pcmData))
	fileSize := dataSize + 36
	io.WriteString(buf, "RIFF")
	binary.Write(buf, binary.LittleEndian, fileSize)
	io.WriteString(buf, "WAVE")

	io.WriteString(buf, "fmt ")
	binary.Write(buf, binary.LittleEndian, uint32(16)) // Subchunk1Size
	binary.Write(buf, binary.LittleEndian, uint16(1))  // AudioFormat (PCM)
	binary.Write(buf, binary.LittleEndian, channels)
	binary.Write(buf, binary.LittleEndian, sampleRate)

	byteRate := sampleRate * uint32(channels) * 2 // 16-bit = 2 bytes
	binary.Write(buf, binary.LittleEndian, byteRate)

	blockAlign := channels * 2
	binary.Write(buf, binary.LittleEndian, blockAlign)
	binary.Write(buf, binary.LittleEndian, uint16(16)) // BitsPerSample
	io.WriteString(buf, "data")
	binary.Write(buf, binary.LittleEndian, dataSize)

	buf.Write(pcmData)
	return buf.Bytes(), nil
}
