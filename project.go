package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed runner.clj
var PROJECT_SCRIPT string
var FILE_EXTENSIONS = []string{"clj", "cljs", "joke", "edn"}

type Project struct {
	dir       string
	project   string
	extension string
}

func GetProject() *Project {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for _, extension := range FILE_EXTENSIONS {
		path := filepath.Join(currentDir, fmt.Sprintf("project.%s", extension))
		if _, err := os.Stat(path); err == nil {
			return &Project{
				dir:       currentDir,
				project:   path,
				extension: extension,
			}
		}
	}

	log.Fatal(errors.New("Project file not found."))
	return nil
}

func (self Project) GetRunnerCode() string {
	replacer := strings.NewReplacer(
		"{{ FILE_EXTENSION }}", self.extension,
		"{{ PROJECT_FILE }}", self.project,
	)

	return replacer.Replace(PROJECT_SCRIPT)
}

var project *Project
