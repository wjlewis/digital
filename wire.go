package digital

// Wire represents a wire connecting two components in a digital circuit.
type Wire struct {
	Value bool
	ee    *emitter[bool]
}

func NewWire() *Wire {
	return &Wire{Value: false, ee: newEmitter[bool]()}
}

// Set the value of the wire.
func (w *Wire) Set(value bool) {
	if value != w.Value {
		w.Value = value
		w.ee.emit(w.Value)
	}
}

// Listen for changes to the value of this wire.
//
// The listener is also called once when invoking OnUpdate in order to
// initialize downstream components.
func (w *Wire) OnUpdate(listener func(bool)) {
	// This initial invocation of listener is necessary to properly initialize
	// circuit elements.
	listener(w.Value)
	w.ee.on(listener)
}
