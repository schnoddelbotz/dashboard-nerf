package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
	if duration > 0 {
		addDurationLabel(thumb, duration)
	}
	return
}

func addDurationLabel(thumb string, duration float64) {
	reader, err := os.Open(thumb)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	img, _ := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	durationString := fmt.Sprintf("%.2fs", duration)
	boxLength := (len(durationString) * 8) + 2 // y = 8x + 2 -> enough pixels per character
	mask := image.NewRGBA(image.Rect(0, 0, boxLength, 15))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(img.(draw.Image), mask.Bounds(), &image.Uniform{black}, image.ZP, draw.Src)

	white := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(400), Y: fixed.Int26_6(700)}
	d := &font.Drawer{
		Dst:  img.(draw.Image),
		Src:  image.NewUniform(white),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(durationString)
	f, err := os.Create(thumb)

	err = png.Encode(f, img)
	f.Close()
}

func getVideoDuration(video string) (duration float64, err error) {
	var args = []string{"-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1"}
	cmd := exec.Command("ffprobe", append(args, video)...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return
	}
	duration, err = strconv.ParseFloat(strings.TrimRight(string(output), "\n"), 32)
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
