package digital

import (
	"fmt"
)

// Circuit represents a digital circuit composed on NAND gates, inputs, and
// outputs.
type Circuit struct {
	pendingEvents []event
	inputs        map[string]*Input
	outputs       map[string]*Output
	step          int64
	stepMsgs      []string
	stepEmitter   *emitter[int64]
}

type event struct {
	wire  *Wire
	value bool
}

// Construct a new, empty circuit.
func NewCircuit() *Circuit {
	return &Circuit{
		step:          0,
		pendingEvents: make([]event, 0),
		inputs:        make(map[string]*Input),
		outputs:       make(map[string]*Output),
		stepMsgs:      make([]string, 0),
		stepEmitter:   newEmitter[int64](),
	}
}

// Add a NAND gate to the circuit, with inputs in1 and in2, and output out.
func (ct *Circuit) AddNand(in1, in2, out *Wire) {
	handleUpdate := func(_ bool) {
		value := !(in1.Value && in2.Value)
		ct.addPendingEvent(out, value)
	}
	in1.OnUpdate(handleUpdate)
	in2.OnUpdate(handleUpdate)
}

// Connect an input with the provided name.
func (ct *Circuit) ConnectInput(input *Input, name string) {
	ct.inputs[name] = input
}

// Connect an output with the provided name.
func (ct *Circuit) ConnectOutput(output *Output, name string) {
	ct.outputs[name] = output
}

// Call listener each time a step occurs, passing the current step.
func (ct *Circuit) OnStep(listener func(int64)) {
	ct.stepEmitter.on(listener)
}

// Simulate the circuit. Steps occur ever stepInterval milliseconds, and logs
// are written to logFilename.txt.
func (ct *Circuit) Run(stepInterval int, logFilename string) {
	ct.runWith(newCliRunner(stepInterval, logFilename))
}

func (ct *Circuit) addPendingEvent(wire *Wire, value bool) {
	evt := event{wire, value}
	ct.pendingEvents = append(ct.pendingEvents, evt)
}

func (ct *Circuit) addStepMsg(msg string) {
	ct.stepMsgs = append(ct.stepMsgs, msg)
}

func (ct *Circuit) clearStepMsgs() {
	ct.stepMsgs = make([]string, 0)
}

func (ct *Circuit) runWith(r runner) {
	for name, output := range ct.outputs {
		name := name
		output := output
		output.onUpdate(func(_ struct{}) {
			value, err := output.value()
			if err != nil {
				msg := fmt.Sprintf("error displaying %s's value: %s", name, err.Error())
				ct.addStepMsg(msg)
			} else {
				msg := fmt.Sprintf("%s -> %s", name, value)
				ct.addStepMsg(msg)
			}
		})
	}

	defer r.close()
	go r.start()

	for {
		select {
		case <-r.steps():
			ct.runStep(r)
		case cmd := <-r.cmds():
			shouldExit := ct.runCmd(cmd, r)
			if shouldExit {
				return
			}
		}
	}
}

func (ct *Circuit) runStep(r runner) {
	ct.step += 1
	ct.stepEmitter.emit(ct.step)

	events := ct.pendingEvents
	ct.pendingEvents = make([]event, 0)

	for _, event := range events {
		event.wire.Set(event.value)
	}

	for _, msg := range ct.stepMsgs {
		r.log(fmt.Sprintf("%d: %s", ct.step, msg))
	}
	ct.clearStepMsgs()
}

func (ct *Circuit) runCmd(text string, r runner) bool {
	cmd, err := parseCmd(text)
	if err != nil {
		r.respond(err.Error())
		return false
	}

	switch c := cmd.(type) {
	case getCmd:
		ct.runGetCmd(c, r)
		return false
	case setCmd:
		ct.runSetCmd(c, r)
		return false
	case exitCmd:
		return true
	default:
		panic("bad command")
	}
}

func (ct *Circuit) runGetCmd(cmd getCmd, r runner) {
	output := ct.outputs[cmd.name]
	if output == nil {
		r.respond("output doesn't exist")
		return
	}

	value, err := output.value()
	if err != nil {
		r.respond(err.Error())
		return
	}

	r.respond(value)
}

func (ct *Circuit) runSetCmd(cmd setCmd, r runner) {
	input := ct.inputs[cmd.name]
	if input == nil {
		r.respond("input doesn't exist")
		return
	}

	values, err := input.parse(cmd.value)
	if err != nil {
		r.respond(err.Error())
		return
	}

	wantCount := len(input.wires)
	gotCount := len(values)
	if gotCount != wantCount {
		r.respond(fmt.Sprintf("incorrect number of values: want %d, got %d", wantCount, gotCount))
		return
	}

	for i, wire := range input.wires {
		wire.Set(values[i])
	}
	ct.addStepMsg(fmt.Sprintf("%s set to %s", cmd.name, cmd.value))
	r.respond("set!")
}
