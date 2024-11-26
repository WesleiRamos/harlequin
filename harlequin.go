package main

import (
	"flag"
	"os"
)

func main() {
	watch := flag.Bool("watch", false, "Watch files changes")
	flag.Parse()

	project = GetProject()
	runner = CreateRunner(project.GetRunnerCode())
	defer os.Remove(runner.file.Name())

	if *watch {
		WatchFiles()
	} else {
		runner.Run()
	}
}
