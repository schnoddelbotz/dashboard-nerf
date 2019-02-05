package main

import (
	"os/exec"
	"path/filepath"
	"strings"
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
	for elem := range playQueue {
		switch elem.MediaType {
		case audio:
			args := audioPlayer[1:]
			args = append(args, elem.Filename)
			cmd = exec.Command(audioPlayer[0], args...)
		case video:
			args := videoPlayer[1:]
			args = append(args, elem.Filename)
			cmd = exec.Command(videoPlayer[0], args...)
		case speech:
			args := speechCommand[1:]
			args = append(args, elem.TextToSpeak)
			cmd = exec.Command(speechCommand[0], args...)
		}
		cmd.Run()
	}
}

func getMediaType(filename string) mediaType {
	switch extension := strings.ToLower(filepath.Ext(filename)); extension {
	case ".mp3":
		fallthrough
	case ".wav":
		return audio
	}
	return video
}
