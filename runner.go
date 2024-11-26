package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Runner struct {
	cmd  *exec.Cmd
	file *os.File
	args []string
}

func CreateRunner(code string) *Runner {
	file, err := os.CreateTemp("", "runner.clj")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := file.Write([]byte(code)); err != nil {
		log.Fatal(err)
	}

	return &Runner{file: file, args: append([]string{file.Name()}, flag.Args()...)}
}

func (self *Runner) Run() {
	self.cmd = exec.Command("joker", self.args...)
	self.cmd.Stdout = os.Stdout
	self.cmd.Stderr = os.Stderr
	if err := self.cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "\nFailed to run:", err)
	}
}

func (self *Runner) Kill() {
	if self.cmd.Process != nil {
		self.cmd.Process.Kill()
	}
}

var runner *Runner