package gotask

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

const (
	// the task is received by a worker
	StateReceived = "RECEIVED"
	// the task is started
	StateStarted = "STARTED"
	// the task is processed successfully
	StateSuccess = "SUCCESS"
	// the processing of the task fails
	StateFailure = "FAILURE"
)

type Result interface{}

type State struct {
	UUID   string
	State  string
	Result Result
	Error  error
}

type Task interface {
	Run() (Result, error)
	Notify(state *State) error
}

type Signature struct {
	UUID string
	Name string
	Args map[string]interface{}
}

type Constructor func() Task

func safelyRun(task Task) (result Result, err error) {
	defer func() {
		// Recover from panic and set err.
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				err = errors.New("Running task caused a panic")
			}
		}
	}()

	result, err = task.Run()
	return
}

func Process(registry map[string]Constructor, sig *Signature) error {
	constructor, ok := registry[sig.Name]
	if !ok {
		return errors.New("No task named " + sig.Name)
	}

	task := constructor()

	state := &State{UUID: sig.UUID, State: StateReceived}
	if err := task.Notify(state); err != nil {
		return err
	}

	// Decode map[string]interface{} into a specific task struct
	if err := mapstructure.Decode(sig.Args, task); err != nil {
		return err
	}

	state.State = StateStarted
	if err := task.Notify(state); err != nil {
		return err
	}

	result, err := safelyRun(task)
	if err != nil {
		state.State = StateFailure
		state.Error = err
	} else {
		state.State = StateSuccess
		state.Result = result
	}
	err = task.Notify(state)
	return err
}
