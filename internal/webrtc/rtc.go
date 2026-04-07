package webrtc

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"project_for_tmk_04_06/internal/ai/tts"
	"project_for_tmk_04_06/internal/audio"
	"project_for_tmk_04_06/internal/runner"
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

// SDP to base64  / base64 to SDP

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
	// 定义 ICEConnectionStateChange 的钩子
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		pterm.Info.Printf("RTC ICE Connection State has changed: %s\n", connectionState.String())
	})
}

func (r *RTCManager) setupDataChannel(dc *webrtc.DataChannel) {
	// Open 的钩子
	dc.OnOpen(func() {
		pterm.Success.Printf("Data channel '%s'-'%d' open. Ready to communicate!\n", dc.Label(), dc.ID())
	})

	// Message 的钩子
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
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	con, err := r.api.NewPeerConnection(config)
	if err != nil {
		return err
	}

	defer con.Close()

	r.setupCallbacks(con)
	dc, err := con.CreateDataChannel("translations", nil)
	if err != nil {
		return err
	}
	r.setupDataChannel(dc)
	// 钩子的设置

	offer, err := con.CreateOffer(nil) // Host 的自我介绍信
	if err != nil {
		return err
	}
	if err = con.SetLocalDescription(offer); err != nil {
		return err
	}
	<-webrtc.GatheringCompletePromise(con)

	encodedOffer, _ := Encode(con.LocalDescription())

	pterm.DefaultHeader.Println("🏠 RTC Host Mode")
	pterm.Warning.Println("Please copy the line below and share it with the joining peer as their <room-id>:")

	fmt.Printf("\n%s\n\n", encodedOffer)

	pterm.Info.Println("Paste the Answer from joining peer below and press ENTER:")

	answerStr := readStdin()

	var answer webrtc.SessionDescription
	if err := Decode(answerStr, &answer); err != nil {
		return fmt.Errorf("failed to decode answer: %w", err)
	}
	if err := con.SetRemoteDescription(answer); err != nil {
		return fmt.Errorf("failed to set remote description: %w", err)
	}

	// 连接成功
	streamRn := runner.NewStreamRunner()

	streamRn.Callback = func(source, translation string) {
		payload := Payload{
			Language: targetLang,
			Text:     translation,
		}
		if data, err := json.Marshal(payload); err == nil {
			_ = dc.Send(data)
		}
	}

	return streamRn.Run(ctx, sourceLang, targetLang, ttsEnabled)
}

func (r *RTCManager) Join(ctx context.Context, roomID string) error {
	return nil
}
