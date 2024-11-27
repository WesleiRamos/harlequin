package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var FILE_EXTENSIONS = []string{"clj", "cljs", "joke", "edn"}

//go:embed runner.clj
var PROJECT_RUNNER string

//go:embed template/project.joke
var PROJECT_SCRIPT string

//go:embed template/src/project/main.joke
var PROJECT_ENTRYPOINT string

type Project struct {
	dir       string
	project   string
	extension string
}

func NewProject() *Project {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Project{dir: currentDir}
}

func (self *Project) GetCurrentProject() {
	for _, extension := range FILE_EXTENSIONS {
		path := filepath.Join(self.dir, fmt.Sprintf("project.%s", extension))
		if _, err := os.Stat(path); err == nil {
			self.project = path
			self.extension = extension
			return
		}
	}

	log.Fatal(errors.New("Project file not found."))
}

func (self Project) GetRunnerCode() string {
	replacer := strings.NewReplacer(
		"{{ FILE_EXTENSION }}", self.extension,
		"{{ PROJECT_FILE }}", self.project,
	)

	return replacer.Replace(PROJECT_RUNNER)
}

func (self Project) NewProject(name string) {
	fmt.Printf("\n > harlequin new %s\n\n", name)

	if _, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name); err != nil {
		log.Fatal(errors.New("Invalid project name."))
	}

	folder := fmt.Sprintf("%s/src/%s", name, name)
	if err := os.MkdirAll(folder, 0750); err != nil {
		log.Fatal(err)
	}

	replacer := strings.NewReplacer("{{ PROJECT_NAME }}", name)
	project := replacer.Replace(PROJECT_SCRIPT)

	if err := os.WriteFile(fmt.Sprintf("%s/project.joke", name), []byte(project), 0644); err != nil {
		log.Fatal(err)
	}

	entrypoint := replacer.Replace(PROJECT_ENTRYPOINT)
	if err := os.WriteFile(fmt.Sprintf("%s/main.joke", folder), []byte(entrypoint), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\033[32mSuccess! \033[0m")
	fmt.Printf("Created '%s' at '%s/%s'\n", name, self.dir, name)
	fmt.Printf("Access it with '\033[33mcd\033[0m %s' and run with \033[33mharlequin\033[0m\n", name)
}

var project *Project
