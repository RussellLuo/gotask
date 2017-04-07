package tasks

import (
	"errors"
	"log"

	"github.com/RussellLuo/gotask"
)

type Add struct {
	X, Y int
}

func (a *Add) Run() (gotask.Result, error) {
	return a.X + a.Y, nil
}

func (a *Add) Notify(state *gotask.State) error {
	log.Printf("state: %#v", state)
	return nil
}

type Greet struct {
	Words string
}

func (g *Greet) Run() (gotask.Result, error) {
	return "Hello, " + g.Words, nil
}

func (g *Greet) Notify(state *gotask.State) error {
	log.Printf("state: %#v", state)
	return nil
}

type Panic struct{}

func (p *Panic) Run() (gotask.Result, error) {
	panic(errors.New("oops"))
}

func (p *Panic) Notify(state *gotask.State) error {
	log.Printf("state: %#v", state)
	return nil
}
