package main

import (
	"fmt"
	"log"
	"os/exec"
)

type mediaType int

type playRequest struct {
	MediaType   mediaType
	Filename    string
	TextToSpeak string
}

const (
	audio  mediaType = iota
	video  mediaType = iota
	speech mediaType = iota
)

var (
	playQueue = make(chan playRequest)
	cmd       *exec.Cmd
)

func startQueuePlayer() {
	var runCmd string
	for elem := range playQueue {
		switch elem.MediaType {
		case audio:
			playAudio(elem.Filename)
		case video:
			runCmd = fmt.Sprintf(videoPlayer, elem.Filename)
			log.Printf("Executing: %s\n", runCmd)
			cmd = exec.Command("sh", "-c", runCmd)
			cmd.Run()
		case speech:
			runCmd = fmt.Sprintf(speechCommand, elem.TextToSpeak)
			log.Printf("Executing: %s\n", runCmd)
			cmd = exec.Command("sh", "-c", runCmd)
			cmd.Run()
		}
	}
}
