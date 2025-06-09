package digital

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type cliRunner struct {
	tickInterval  int
	ticksChan     chan struct{}
	cmdsChan      chan string
	responsesChan chan string
	logFile       *os.File
	tty           *bufio.Scanner
}

func newCliRunner(tickInterval int, logFilename string) *cliRunner {
	logFile, err := os.Create(logFilename)
	if err != nil {
		panic(err)
	}

	tty := bufio.NewScanner(os.Stdin)

	ticksChan := make(chan struct{})
	cmdsChan := make(chan string)
	responsesChan := make(chan string)

	return &cliRunner{
		tickInterval,
		ticksChan,
		cmdsChan,
		responsesChan,
		logFile,
		tty,
	}
}

func (r *cliRunner) steps() <-chan struct{} {
	return r.ticksChan
}

func (r *cliRunner) cmds() <-chan string {
	return r.cmdsChan
}

func (r *cliRunner) respond(msg string) {
	r.responsesChan <- msg
}

func (r *cliRunner) log(msg string) {
	r.logFile.WriteString(msg)
	r.logFile.WriteString("\n")
}

func (r *cliRunner) start() {
	go func() {
		for {
			r.ticksChan <- struct{}{}
			time.Sleep(time.Duration(r.tickInterval) * time.Millisecond)
		}
	}()

	for {
		fmt.Print("> ")
		if !r.tty.Scan() {
			continue
		}

		text := strings.TrimSpace(r.tty.Text())
		if len(text) == 0 {
			continue
		}

		r.cmdsChan <- text
		response := <-r.responsesChan
		fmt.Println(response)
	}
}

func (r *cliRunner) close() {
	r.logFile.Close()
}
