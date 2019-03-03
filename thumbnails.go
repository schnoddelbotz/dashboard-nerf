package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"path"
)

func buildThumbnails() {
	if _, err := os.Stat(thumbRoot); os.IsNotExist(err) {
		err := os.Mkdir(thumbRoot, os.ModePerm)
		if err != nil {
			fmt.Printf("FATAL: Cannot create thumb directory %s: %s\n", thumbRoot, err)
			return
		}
	}
	for _, v := range getContent().Videos {
		video := mediaRoot + "/" + v.Filename
		thumbnail := thumbRoot +  "/" +v.Filename + ".png"
		if _, err := os.Stat(thumbnail); os.IsNotExist(err) {
			err = buildThumb(video, thumbnail)
			if err != nil {
				fmt.Printf("Thumb generation failed for %s: %s\n", thumbnail, err)
			}
		}
	}
}

func buildThumb(video string, thumb string) (err error) {
	fmt.Printf("Creating thumb: %s -> %s\n", video, thumb)
	duration, err := getVideoDuration(video)
	if err != nil {
		return
	}
	var args []string
	args = append(args, "-i", video)
	args = append(args, "-vframes", "1")
	args = append(args, "-ss", fmt.Sprintf("%f", duration/2))
	args = append(args, "-vf", "scale=-2:120")
	args = append(args, thumb)
	cmd := exec.Command("ffmpeg", args...)
	err = cmd.Run()
	return
}

func getVideoDuration(video string) (duration float64, err error) {
	var args = []string{"-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1"}
	cmd := exec.Command("ffprobe", append(args,video)...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return
	}
	duration, err = strconv.ParseFloat(strings.TrimRight(string(output),"\n") , 32)
	return
}

func watchForChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					time.Sleep(10 * time.Second) // ugly ... wait, as file might still grow
					thumb := thumbRoot +"/"+ path.Base(event.Name) + ".png"
					if _, err := os.Stat(thumb); os.IsNotExist(err) {
						err := buildThumb(event.Name, thumb)
						if err != nil {
							fmt.Printf("Thumb generation failed for %s: %s\n", event.Name, err)
						}
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					thumb := thumbRoot +"/"+ path.Base(event.Name) + ".png"
					fmt.Printf("Video %s was removed, deleting thumb %s\n", event.Name, thumb)
					if _, err := os.Stat(thumb); err == nil {
						err := os.Remove(thumb)
						if err != nil {
							fmt.Printf("Failed to delete thumb: %s\n", thumb)
						}
					}
				}
			}
		}
	}()

	err = watcher.Add(mediaRoot)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}