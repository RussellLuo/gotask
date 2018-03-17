package gotask

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

const (
	// The task is received by a worker.
	StateReceived = "RECEIVED"
	// The task is started.
	StateStarted = "STARTED"
	// The task is processed successfully.
	StateSuccess = "SUCCESS"
	// The processing of the task fails.
	StateFailure = "FAILURE"
)

// Result is the output of a (successful) task.
type Result interface{}

// State is the processing state of a task.
type State struct {
	ID     string
	State  string
	Result Result
	Error  error
}

// Task represents a job need to be done.
type Task interface {
	Run() (Result, error)
	Notify(state *State) error
}

// Signature is the description of a task.
type Signature struct {
	ID   string
	Name string
	Args map[string]interface{}
}

// Registry is a mapping from task name to task.
type Registry map[string]func() Task

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
				err = errors.New("running task caused a panic")
			}
		}
	}()

	result, err = task.Run()
	return
}

// Process parses out the task described by sig, and then perform the task.
// Any error occurred in the processing will be returned immediately.
func Process(registry Registry, sig *Signature) error {
	constructor, ok := registry[sig.Name]
	if !ok {
		return errors.New("no task named " + sig.Name)
	}

	task := constructor()

	state := &State{ID: sig.ID, State: StateReceived}
	if err := task.Notify(state); err != nil {
		return err
	}

	// Decode map[string]interface{} into a specific task struct.
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
