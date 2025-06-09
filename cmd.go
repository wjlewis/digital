package digital

import (
	"errors"
	"strings"
)

type cmd interface{}

type getCmd struct {
	name string
}

type setCmd struct {
	name  string
	value string
}

type exitCmd struct{}

var errBadCmd = errors.New("bad command")

func parseCmd(text string) (cmd, error) {
	parts := strings.Fields(text)

	if len(parts) == 1 && strings.ToLower(parts[0]) == "exit" {
		return exitCmd{}, nil
	} else if len(parts) == 2 && strings.ToLower(parts[0]) == "get" {
		return getCmd{name: parts[1]}, nil
	} else if len(parts) >= 3 && strings.ToLower(parts[0]) == "set" {
		rest := strings.Join(parts[2:], " ")
		return setCmd{name: parts[1], value: rest}, nil
	} else {
		return nil, errBadCmd
	}
}
