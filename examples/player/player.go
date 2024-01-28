package main

import (
	"fmt"
	"github.com/ebitengine/oto/v3"
	"os"
	"time"

	"github.com/jonchammer/audio-io/core"
	"github.com/jonchammer/audio-io/wave"
)

func main() {

	file, err := os.Open("../../voice.wav")
	if err != nil {
		failF(err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Create a reader using the file
	reader := wave.NewReader(file)
	header, err := reader.Header()
	if err != nil {
		failF(err)
	}

	sampleType, err := header.SampleType()
	if err != nil {
		failF(err)
	}

	// Create a decoder that will translate samples in real time from their
	// native format to Int16 so oto can understand them.
	decoder := core.NewDecoder(reader, sampleType, core.SampleTypeInt16)

	// Set up the player
	otoCtx, readyChan, err := oto.NewContext(
		&oto.NewContextOptions{
			SampleRate:   int(header.FrameRate()),
			ChannelCount: int(header.ChannelCount()),
			Format:       oto.FormatSignedInt16LE,
		},
	)
	if err != nil {
		failF(err)
	}
	<-readyChan
	player := otoCtx.NewPlayer(decoder)
	player.Play()

	// Wait for the sound to finish playing
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	err = player.Close()
	if err != nil {
		failF(err)
	}
}

func failF(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
