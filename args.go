package main

import (
	"errors"
	"fmt"
	"os"
)

type Arg struct {
	name  string
	index int
}

func GetArg(name string) *Arg {
	for index, arg := range os.Args {
		if arg == name {
			return &Arg{name, index}
		}
	}

	return &Arg{name, -1}
}

func (self Arg) Value(pos int) (string, error) {
	if self.index+pos >= len(os.Args) {
		return "", errors.New(fmt.Sprintf("No value for %s", self.name))
	}

	return os.Args[self.index+pos], nil
}
