package main

// This code is slightly extended beep example code from
// https://github.com/faiface/beep/blob/master/examples/playing/mp3-playing.go

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func playAudio(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening %s: %s", filename, err)
		return
	}

	var s beep.StreamSeekCloser
	var format beep.Format

	switch extension := strings.ToLower(filepath.Ext(filename)); extension {
	case ".wav":
		s, format, err = wav.Decode(f)
	case ".mp3":
		s, format, err = mp3.Decode(f)
	default:
		log.Printf("Unsupported file format: %s", filename)
		return
	}
	if err != nil {
		log.Printf("Audio playback error: %s", err)
		return
	}

	// Init the Speaker with the SampleRate of the format and a buffer size of 1/10s
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// Channel, which will signal the end of the playback.
	playing := make(chan struct{})

	// Now we Play our Streamer on the Speaker
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		// Callback after the stream Ends
		close(playing)
	})))
	<-playing
}
