package digital

import "errors"

// Input represents an input device that can be used to set the values of one or
// more wires in a circuit.
type Input struct {
	wires []*Wire
	parse func(string) ([]bool, error)
}

// Construct an input that manipulates the provided wires. The parse function
// dictates how string values are transformed into bit values to be set on the
// wires: it must produce a slice of boolean values whose length is equal to the
// number of wires.
func NewInput(wires []*Wire, parse func(string) ([]bool, error)) *Input {
	return &Input{wires, parse}
}

// A basic input simulating a button that's either pressed (1) or not (0).
func Button(wire *Wire) *Input {
	parse := func(value string) ([]bool, error) {
		return []bool{value != "0"}, nil
	}
	return NewInput([]*Wire{wire}, parse)
}

// Output represents an output device that responds to changes in one or more
// wires in a circuit.
type Output struct {
	wires   []*Wire
	unparse func([]bool) (string, error)
	ee      *emitter[struct{}]
}

// Construct an output the responds to changes in the provided wires. The
// unparse function determines how bit values are presented to the simulator's
// user interface: it must expect a slice of boolean values whose length is
// equal to the number of wires.
func NewOutput(wires []*Wire, unparse func([]bool) (string, error)) *Output {
	ee := newEmitter[struct{}]()
	out := &Output{wires, unparse, ee}
	for _, wire := range wires {
		wire.OnUpdate(func(_ bool) {
			out.ee.emit(struct{}{})
		})
	}
	return out
}

func (o *Output) onUpdate(listener func(struct{})) {
	o.ee.on(listener)
}

func (o *Output) value() (string, error) {
	values := make([]bool, len(o.wires))
	for i, wire := range o.wires {
		values[i] = wire.Value
	}

	return o.unparse(values)
}

// A basic output simulating a bulb that's either on (1) or off (0).
func Bulb(wire *Wire) *Output {
	unparse := func(bits []bool) (string, error) {
		if len(bits) != 1 {
			return "", errors.New("expected exactly 1 value")
		}

		if bits[0] {
			return "1", nil
		} else {
			return "0", nil
		}
	}
	return NewOutput([]*Wire{wire}, unparse)
}
