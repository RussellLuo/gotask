package gotask_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/RussellLuo/gotask"
)

type Base struct {
	NotifyFunc func(state *gotask.State) error
}

func (b *Base) Notify(state *gotask.State) error {
	return b.NotifyFunc(state)
}

type Add struct {
	X, Y int
	Base
}

func (a *Add) Run() (gotask.Result, error) {
	return a.X + a.Y, nil
}

type Greet struct {
	Words string
	Base
}

func (g *Greet) Run() (gotask.Result, error) {
	return "Hello, " + g.Words, nil
}

type Panic struct {
	Base
}

func (p *Panic) Run() (gotask.Result, error) {
	panic(errors.New("oops"))
}

func TestProcess(t *testing.T) {
	type WantType struct {
		err    error
		states []gotask.State
	}
	cases := []struct {
		in   gotask.Signature
		want WantType
	}{
		{
			in: gotask.Signature{
				UUID: "uuid-1", Name: "add",
				Args: map[string]interface{}{"x": 1, "y": 2},
			},
			want: WantType{
				err: nil,
				states: []gotask.State{
					{
						UUID:   "uuid-1",
						State:  "RECEIVED",
						Result: nil,
						Error:  nil,
					},
					{
						UUID:   "uuid-1",
						State:  "STARTED",
						Result: nil,
						Error:  nil,
					},
					{
						UUID:   "uuid-1",
						State:  "SUCCESS",
						Result: 3,
						Error:  nil,
					},
				},
			},
		},
		{
			in: gotask.Signature{
				UUID: "uuid-2", Name: "greet",
				Args: map[string]interface{}{"words": "Russell"},
			},
			want: WantType{
				err: nil,
				states: []gotask.State{
					{
						UUID:   "uuid-2",
						State:  "RECEIVED",
						Result: nil,
						Error:  nil,
					},
					{
						UUID:   "uuid-2",
						State:  "STARTED",
						Result: nil,
						Error:  nil,
					},
					{
						UUID:   "uuid-2",
						State:  "SUCCESS",
						Result: "Hello, Russell",
						Error:  nil,
					},
				},
			},
		},
		{
			in: gotask.Signature{
				UUID: "uuid-3", Name: "panic",
				Args: map[string]interface{}{},
			},
			want: WantType{
				err: nil,
				states: []gotask.State{
					{
						UUID:   "uuid-3",
						State:  "RECEIVED",
						Result: nil,
						Error:  nil,
					},
					{
						UUID:   "uuid-3",
						State:  "STARTED",
						Result: nil,
						Error:  nil,
					},
					{
						UUID:   "uuid-3",
						State:  "FAILURE",
						Result: nil,
						Error:  errors.New("oops"),
					},
				},
			},
		},
		{
			in: gotask.Signature{
				UUID: "uuid-4", Name: "unknown",
				Args: map[string]interface{}{},
			},
			want: WantType{
				err:    errors.New("No task named unknown"),
				states: []gotask.State{},
			},
		},
	}

	makeRegistry := func(notify func(state *gotask.State) error) map[string]gotask.Constructor {
		return map[string]gotask.Constructor{
			"add": func() gotask.Task {
				return &Add{Base: Base{NotifyFunc: notify}}
			},
			"greet": func() gotask.Task {
				return &Greet{Base: Base{NotifyFunc: notify}}
			},
			"panic": func() gotask.Task {
				return &Panic{Base: Base{NotifyFunc: notify}}
			},
		}
	}

	for _, c := range cases {
		states := make([]gotask.State, 0, 3)
		registry := makeRegistry(func(state *gotask.State) error {
			states = append(states, *state)
			return nil
		})
		err := gotask.Process(registry, &c.in)
		if !reflect.DeepEqual(err, c.want.err) {
			t.Errorf("Got (%#v) != Want (%#v)", err, c.want.err)
		}
		if !reflect.DeepEqual(states, c.want.states) {
			t.Errorf("Got (%#v) != Want (%#v)", states, c.want.states)
		}
	}
}
