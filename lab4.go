package main

import (
	"fmt"
	"./engine"
	"bufio"
	"strings"
	"os"
	"errors"
)

type printCommand struct {
	arg string
}

func (p *printCommand) Execute(loop engine.Handler) {
	fmt.Println(p.arg)
}

type catCommand struct {
	arg1, arg2 string
}

func (cat *catCommand) Execute(loop engine.Handler) {
	res := cat.arg1 + cat.arg2
	loop.Post(&printCommand{arg: res})
}

func parse(commandLine string) engine.Command {
	parts := strings.Fields(commandLine)
	if parts[0] == "print" {
		if len(parts) == 2 {
			return &printCommand{arg: parts[1]}
		} else if len(parts) == 1 {
			return &printCommand{arg: "SYNTAX ERROR: " + errors.New("no argument").Error()}
		} else {
			return &printCommand{arg: "SYNTAX ERROR: " + errors.New("too many arguments").Error()}
		}
	} else if parts[0] == "cat" {
		if len(parts) == 3 {
			return &catCommand{arg1: parts[1], arg2: parts[2]}
		} else if len(parts) < 3 {
			return &printCommand{arg: "SYNTAX ERROR: " + errors.New("not enough arguments").Error()}
		} else {
			return &printCommand{arg: "SYNTAX ERROR: " + errors.New("too many arguments").Error()}
		}
	} else {
		return &printCommand{arg: "SYNTAX ERROR: " + errors.New("unexpected command").Error()}
	}
}

func main() {
	eventLoop := new(engine.EventLoop)
	eventLoop.Start()
	if input, err := os.Open("./commands.txt"); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := parse(commandLine) // parse the line to get an instance of Command
			eventLoop.Post(cmd)
		}
	}
	eventLoop.AwaitFinish()
}