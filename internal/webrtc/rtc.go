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
		ttsClient: tts.NewSiliconFlowTTS(),
	}
}

func Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

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
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	pc, err := r.api.NewPeerConnection(config)
	if err != nil {
		return err
	}
	defer pc.Close()

	r.setupCallbacks(pc)

	dc, err := pc.CreateDataChannel("translations", nil)
	if err != nil {
		return err
	}
	r.setupDataChannel(dc)

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		return err
	}
	if err = pc.SetLocalDescription(offer); err != nil {
		return err
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)
	<-gatherComplete

	encodedOffer, _ := Encode(pc.LocalDescription())
	pterm.DefaultHeader.Println("🏠 RTC Host Mode")
	pterm.Warning.Println("Please copy the line below and share it with the joining peer as their <room-id>:")
	fmt.Printf("\n%s\n\n", encodedOffer)

	pterm.Info.Println("Paste the Answer from joining peer below and press ENTER:")
	answerStr := readStdin()

	var answer webrtc.SessionDescription
	if err := Decode(answerStr, &answer); err != nil {
		return fmt.Errorf("failed to decode answer: %w", err)
	}

	if err := pc.SetRemoteDescription(answer); err != nil {
		return fmt.Errorf("failed to set remote description: %w", err)
	}

	// We are connected. Start localized StreamRunner that forwards translation
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
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	pc, err := r.api.NewPeerConnection(config)
	if err != nil {
		return err
	}
	defer pc.Close()

	r.setupCallbacks(pc)

	pc.OnDataChannel(func(d *webrtc.DataChannel) {
		r.setupDataChannel(d)
	})

	var offer webrtc.SessionDescription
	if err := Decode(roomID, &offer); err != nil {
		return fmt.Errorf("invalid room ID: %w", err)
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		return fmt.Errorf("failed to set remote offer: %w", err)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return err
	}
	if err := pc.SetLocalDescription(answer); err != nil {
		return err
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)
	<-gatherComplete

	encodedAnswer, _ := Encode(pc.LocalDescription())
	pterm.DefaultHeader.Println("🤝 RTC Join Mode")
	pterm.Success.Println("Connected to offer! Copy the answer below and provide it to the Host:")
	fmt.Printf("\n%s\n\n", encodedAnswer)

	<-ctx.Done()
	return nil
}
