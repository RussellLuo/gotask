package gotask_test

import (
	"fmt"

	"github.com/RussellLuo/gotask"
)

type Sum struct {
	X, Y int
}

func (s *Sum) Run() (gotask.Result, error) {
	return s.X + s.Y, nil
}

func (s *Sum) Notify(state *gotask.State) error {
	if state.State == gotask.StateSuccess {
		fmt.Printf("sum: %d\n", state.Result.(int))
	}
	return nil
}

func Example() {
	registry := gotask.Registry{
		"sum":   func() gotask.Task { return &Sum{} },
	}

	err := gotask.Process(registry, &gotask.Signature{
		ID: "id",
		Name: "sum",
		Args: map[string]interface{}{
			"X": 20,
			"Y": 76,
		},
	})
	fmt.Printf("err: %v\n", err)

	// Output:
	// sum: 96
	// err: <nil>
}
