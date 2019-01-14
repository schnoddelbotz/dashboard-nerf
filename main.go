package main

import (
	"flag"
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
	flag.StringVar(&speechCommand, "speech", "say", "speech command")
	flag.StringVar(&mediaRoot, "media", "media", "media root directory")
	flag.StringVar(&httpPort, "port", ":9999", "HTTP server TCP port")
	flag.StringVar(&audioPlayer, "audioplayer", "play", "audioplayer to use for MP3s")
	flag.StringVar(&videoPlayer, "videoplayer", "cvlc --fullscreen", "videoplayer to use for MP4s/WEBM")
	flag.Parse()
	go startQueuePlayer()
	runWebserver(mediaRoot, httpPort)
}
