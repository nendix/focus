package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Player struct {
	context *oto.Context
}

func New() (*Player, error) {
	op := &oto.NewContextOptions{
		SampleRate:   48000, // Match the MP3 file sample rate
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	ctx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}

	<-readyChan

	return &Player{context: ctx}, nil
}

func (p *Player) playSound(filename string) error {
	// Get the path to the audio file
	audioPath := filepath.Join("assets", filename)

	file, err := os.Open(audioPath)
	if err != nil {
		return fmt.Errorf("failed to open audio file %s: %w", filename, err)
	}
	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return fmt.Errorf("failed to decode MP3 %s: %w", filename, err)
	}

	// Create a player from the decoder
	player := p.context.NewPlayer(decoder)
	defer player.Close()

	player.Play()

	// Calculate approximate duration and wait for playback
	sampleRate := decoder.SampleRate()
	length := decoder.Length()
	duration := time.Duration(length) * time.Second / time.Duration(sampleRate) / 4 // 4 bytes per sample (16-bit stereo)

	time.Sleep(duration)

	return nil
}

func (p *Player) PlayWorkEndSound() {
	if err := p.playSound("work-end.mp3"); err != nil {
		// Silently fail if audio can't play, don't crash the app
		fmt.Printf("Audio playback failed: %v\n", err)
	}
}

func (p *Player) PlayBreakEndSound() {
	if err := p.playSound("break-end.mp3"); err != nil {
		// Silently fail if audio can't play, don't crash the app
		fmt.Printf("Audio playback failed: %v\n", err)
	}
}

func (p *Player) Close() {
	if p.context != nil {
		p.context.Suspend()
	}
}
