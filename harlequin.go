package main

import (
	"log"
	"os"
)

func main() {
	project = NewProject()

	if new := GetArg("new"); new.index != -1 {
		projectName, err := new.Value(1)
		if err != nil {
			log.Fatal(err)
		}

		project.NewProject(projectName)
		return
	}

	project.GetCurrentProject()
	runner = CreateRunner(project.GetRunnerCode())
	defer os.Remove(runner.file.Name())

	if watch := GetArg("watch"); watch.index != -1 {
		WatchFiles()
		return
	}

	runner.Run()
}
