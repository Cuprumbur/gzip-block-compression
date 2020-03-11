package command

import (
	"bytes"
)

type Block struct {
	Indx int
	R    *bytes.Reader
}

type Command interface {
	Run(jobs <-chan Block) <-chan Block
}

type Context interface {
	GetStrategy(command string) Command
}

type context struct {
	commands map[string]Command
}

func NewContext(commands map[string]Command) Context {
	return context{commands}
}

func (c context) GetStrategy(command string) Command {
	return c.commands[command]
}
