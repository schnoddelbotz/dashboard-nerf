package main

import (
	"flag"
	"fmt"
	"os"
)

type arrayFlags []string

var (
	// AppVersion set at build time
	AppVersion    string
	mediaRoot     string
	httpPort      string
	slackToken    string
	audioPlayer   arrayFlags
	videoPlayer   arrayFlags
	speechCommand arrayFlags
)

func main() {
	flag.Var(&speechCommand, "speech", "speech command")
	flag.StringVar(&mediaRoot, "media", "media", "media root directory")
	flag.StringVar(&httpPort, "port", ":9999", "HTTP server TCP port")
	flag.StringVar(&slackToken, "slackToken", "", "Slack Bot API Token")
	flag.Var(&audioPlayer, "audioplayer", "audio player to use for MP3/WAV")
	flag.Var(&videoPlayer, "videoplayer", "video player to use for MP4/WEBM")
	flag.Parse()
	if _, err := os.Stat(mediaRoot); os.IsNotExist(err) {
		fmt.Printf("Media root directory '%s' (provided via -media argument) does not exist.\n", mediaRoot)
		os.Exit(1)
	}
	if len(videoPlayer) == 0 {
		videoPlayer = defaultVideoPlayer
	}
	if len(audioPlayer) == 0 {
		audioPlayer = defaultAudioPlayer
	}
	if len(speechCommand) == 0 {
		speechCommand = defaultSpeechApp
	}

	if len(slackToken) > 0 {
		go doSlack()
	}

	go startQueuePlayer()

	fmt.Printf("Webserver starting on port %s ...\n", httpPort)
	runWebserver(mediaRoot, httpPort)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arrayFlags) String() string {
	return "unused"
}
