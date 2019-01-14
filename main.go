package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// AppVersion set at build time
	AppVersion    string
	speechCommand string
	mediaRoot     string
	httpPort      string
	audioPlayer   string
	videoPlayer   string
)

func main() {
	flag.StringVar(&speechCommand, "speech", defaultSpeechApp, "speech command")
	flag.StringVar(&mediaRoot, "media", "media", "media root directory")
	flag.StringVar(&httpPort, "port", ":9999", "HTTP server TCP port")
	flag.StringVar(&audioPlayer, "audioplayer", defaultAudioPlayer, "audioplayer to use for MP3/WAV")
	flag.StringVar(&videoPlayer, "videoplayer", defaultVideoPlayer, "videoplayer to use for MP4/WEBM")
	flag.Parse()
	if _, err := os.Stat(mediaRoot); os.IsNotExist(err) {
		fmt.Printf("Media root directory '%s' (provided via -media argument) does not exist.\n", mediaRoot)
		os.Exit(1)
	}
	go startQueuePlayer()
	fmt.Printf("Webserver starting on port %s ...\n", httpPort)
	runWebserver(mediaRoot, httpPort)
}
