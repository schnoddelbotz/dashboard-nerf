package main

import (
	"fmt"
	"os/exec"
)

type mediaType int

type playRequest struct {
	MediaType mediaType
	Filename  string
}

const (
	audio  mediaType = iota
	video  mediaType = iota
	speech mediaType = iota
)

var playQueue = make(chan playRequest)

func startQueuePlayer() {
	var runCmd string
	for elem := range playQueue {
		switch elem.MediaType {
		case audio:
			runCmd = fmt.Sprintf("%s '%s'", audioPlayer, mediaRoot+"/"+elem.Filename)
		case video:
			runCmd = fmt.Sprintf("%s '%s'", videoPlayer, mediaRoot+"/"+elem.Filename)
		case speech:
			runCmd = fmt.Sprintf("%s '%s'", speechCommand, elem.Filename)
		}
		fmt.Printf("Playing: %s type %d using command: %s\n", elem.Filename, elem.MediaType, runCmd)
		cmd := exec.Command("sh", "-c", runCmd)
		cmd.Run()
	}
}
