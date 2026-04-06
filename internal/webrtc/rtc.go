package webrtc

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"project_for_tmk_04_06/internal/ai/tts"
	"project_for_tmk_04_06/internal/audio"
	"strings"

	"github.com/pion/webrtc/v4"
	"github.com/pterm/pterm"
)

// 返回的结果
type Payload struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

type RTCManager struct {
	api       *webrtc.API
	ttsClient tts.TTS
}

func NewRTCManager() *RTCManager {
	return &RTCManager{
		api:       webrtc.NewAPI(),
		ttsClient: tts.NewEdgeTTS(),
	}
}

// 编码
func Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// 解码
func Decode(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

// 定义一下钩子
func (r *RTCManager) setupCallbacks(peerConnection *webrtc.PeerConnection) {
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		pterm.Info.Printf("RTC ICE Connection State has changed: %s\n", connectionState.String())
	})
}
func (r *RTCManager) setupDataChannel(dc *webrtc.DataChannel) {
	dc.OnOpen(func() {
		pterm.Success.Printf("Data channel '%s'-'%d' open. Ready to communicate!\n", dc.Label(), dc.ID())
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		pterm.Info.Println("<<< Received DataChannel Message")
		var payload Payload

		if err := json.Unmarshal(msg.Data, &payload); err == nil {
			pterm.Success.Printf("PEER [%s]: %s\n", payload.Language, payload.Text)
			if payload.Text != "" {
				audioData, err := r.ttsClient.Synthesize(context.Background(), payload.Text, payload.Language)
				if err == nil {
					_ = audio.PlayMP3(audioData)
				}
			}
		}
	})
}

func (r *RTCManager) Host(ctx context.Context, sourceLang, targetLang string, ttsEnabled bool) error {

	return nil
}

func (r *RTCManager) Join(ctx context.Context, roomID string) error {
	return nil
}
