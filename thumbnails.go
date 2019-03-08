package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
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
		thumbnail := thumbRoot + "/" + v.Filename + ".png"
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
	duration, width, height, err := getVideoMetaData(video)
	if err != nil {
		return
	}
	var args []string
	const scaleOption = "scale=-2:120"
	args = append(args, "-i", video)
	args = append(args, "-vframes", "1")
	args = append(args, "-ss", fmt.Sprintf("%f", duration/2))
	if duration > 0 {
		// calc font size relative to the diagonal
		fontSize := math.Sqrt(width*width+height*height) / 2203.0 * 180
		drawTextOptions := ": fontcolor=white: fontsize=" + fmt.Sprintf("%v", int(fontSize)) + ": x=10: y=h-th-10"
		args = append(args, "-filter_complex", "drawtext=text='"+fmt.Sprintf("%.2f", duration)+
			"s'"+drawTextOptions+","+scaleOption)
	} else {
		args = append(args, "-vf", scaleOption)
	}
	args = append(args, thumb)
	cmd := exec.Command("ffmpeg", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		println("error creating thumbnail with 'ffmpeg ", strings.Join(args, " "), "'")
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	return
}

func getVideoMetaData(video string) (duration, width, height float64, err error) {
	var args = []string{"-v", "error", "-show_entries", "stream=duration,width,height", "-of", "default=noprint_wrappers=1:nokey=1"}
	cmd := exec.Command("ffprobe", append(args, video)...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return
	}
	rawValues := strings.Split(string(output), "\n")
	duration, err = strconv.ParseFloat(strings.TrimRight(rawValues[2], "\n"), 32)
	width, err = strconv.ParseFloat(strings.TrimRight(rawValues[0], "\n"), 32)
	height, err = strconv.ParseFloat(strings.TrimRight(rawValues[1], "\n"), 32)
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
					thumb := thumbRoot + "/" + path.Base(event.Name) + ".png"
					if _, err := os.Stat(thumb); os.IsNotExist(err) {
						err := buildThumb(event.Name, thumb)
						if err != nil {
							fmt.Printf("Thumb generation failed for %s: %s\n", event.Name, err)
						}
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					thumb := thumbRoot + "/" + path.Base(event.Name) + ".png"
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
