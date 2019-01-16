package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// tbd: add youtube handler

// Content are HTML template markers
type Content struct {
	Videos  []Video
	Songs   []Song
	Version string
}

// Song is basically a .mp3 or .wav file
type Song struct {
	Filename string
}

// Video is basically a .mp4 or .webm file
type Video struct {
	Filename string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templateBinary, _ := Asset("index.tpl.html")
	tpl, err := template.New("index").Parse(string(templateBinary))
	if err != nil {
		log.Fatalf("Template parsing error: %v\n", err)
	}
	err = tpl.Execute(w, getContent())
	if err != nil {
		log.Fatalf("Template execution error: %v\n", err)
	}
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	requestPath := r.URL.Path[1:]
	css, err := Asset(requestPath)
	if err == nil {
		w.Write(css)
		return
	}
	http.NotFound(w, r)
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	requestPath := r.URL.Path[1:]
	js, err := Asset(requestPath)
	if err == nil {
		w.Write(js)
		return
	}
	http.NotFound(w, r)
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	// maybe... https://github.com/hajimehoshi/go-mp3/blob/master/example/main.go ?
	filename := strings.TrimPrefix(r.URL.Path[1:], "play/audio/")
	filename = mediaRoot + "/" + filename
	if _, err := os.Stat(filename); err == nil {
		playQueue <- playRequest{Filename: filename, MediaType: audio}
		return
	}
	http.NotFound(w, r)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	// maybe merge with audioHandler -> playHandler ....
	filename := strings.TrimPrefix(r.URL.Path[1:], "play/video/")
	filename = mediaRoot + "/" + filename
	if _, err := os.Stat(filename); err == nil {
		playQueue <- playRequest{Filename: filename, MediaType: video}
		return
	}
	http.NotFound(w, r)
}

func speechHandler(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	reg := regexp.MustCompile("[^a-zA-Z0-9 \\?!\\.,]+")
	text = reg.ReplaceAllString(text, "")
	playQueue <- playRequest{TextToSpeak: text, MediaType: speech}
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	if cmd != nil {
		cmd.Process.Kill()
	}
}

func getContent() Content {
	var songs []Song
	var videos []Video

	files, err := ioutil.ReadDir(mediaRoot)
	if err != nil {
		log.Fatalf("Error reading media root directory: %s", err)
	}

	for _, file := range files {
		filename := file.Name()
		switch extension := strings.ToLower(filepath.Ext(filename)); extension {
		case ".mp3":
			fallthrough
		case ".wav":
			songs = append(songs, Song{Filename: filepath.Base(filename)})
		case ".mp4":
			fallthrough
		case ".webm":
			videos = append(videos, Video{Filename: filepath.Base(filename)})
		}
	}

	return Content{Songs: songs, Videos: videos, Version: AppVersion}
}

func runWebserver(documentRoot string, port string) {
	fs := http.FileServer(http.Dir(documentRoot))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/speech/", speechHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/stop", stopHandler)
	http.HandleFunc("/css/", cssHandler)
	http.HandleFunc("/play/audio/", audioHandler)
	http.HandleFunc("/play/video/", videoHandler)
	http.HandleFunc("/js/", jsHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
