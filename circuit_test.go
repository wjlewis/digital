package digital

import (
	"slices"
	"sync"
	"testing"
)

func TestCircuit(t *testing.T) {
	ct := NewCircuit()

	a := NewWire()
	b := NewWire()
	out := NewWire()

	ct.AddNand(a, b, out)

	ct.ConnectInput(Button(a), "a")
	ct.ConnectInput(Button(b), "b")
	ct.ConnectOutput(Bulb(out), "out")

	r := newTestRunner()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ct.runWith(r)
	}()

	r.step()
	r.sendCmd("set a 1")
	r.sendCmd("set b 1")
	r.step()
	r.sendCmd("exit")

	wg.Wait()

	expected := []string{
		"1: out -> 1",
		"2: a set to 1",
		"2: b set to 1",
		"2: out -> 0",
	}
	if !slices.Equal(r.logData, expected) {
		t.Errorf("expected %q, got %q", expected, r.logData)
	}
}
