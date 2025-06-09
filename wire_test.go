package digital

import (
	"fmt"
)

func ExampleWire_OnUpdate() {
	w := NewWire()

	w.OnUpdate(func(_ bool) {
		fmt.Printf("set to %v\n", w.Value)
	})
	w.Set(true)
	w.Set(true)
	// Output:
	// set to false
	// set to true
}
