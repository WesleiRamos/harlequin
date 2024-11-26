package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bep/debounce"
	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func hasExtension(path, extension string) bool {
	ext := filepath.Ext(path)
	return ext == fmt.Sprintf(".%s", extension)
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func watchFilesChange() {
	debounced := debounce.New(100 * time.Millisecond)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// only watch create and write events
			if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) {
				continue
			}

			// prevent watching files that are not the project extension
			if !hasExtension(event.Name, project.extension) {
				continue
			}

			debounced(func() {
				log.Print("\033[H\033[2J")
				runner.Kill()
				runner.Run()
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func WatchFiles() {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(project.dir, watchDir); err != nil {
		log.Fatal(err)
	}

	go watchFilesChange()
	<-make(chan struct{})
}
