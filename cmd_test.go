package digital

import (
	"testing"
)

type testCase struct {
	text string
	cmd  cmd
	err  error
}

func TestParseCmd(t *testing.T) {
	tests := []testCase{
		{text: "", cmd: nil, err: errBadCmd},
		{text: "get a", cmd: getCmd{name: "a"}, err: nil},
		{text: "set x 1011", cmd: setCmd{name: "x", value: "1011"}, err: nil},
		{text: "set y", cmd: nil, err: errBadCmd},
		{text: "  get  foo ", cmd: getCmd{name: "foo"}, err: nil},
		{text: "GET Bar", cmd: getCmd{name: "Bar"}, err: nil},
		{text: "SeT x HELLO", cmd: setCmd{name: "x", value: "HELLO"}, err: nil},
	}

	for _, tt := range tests {
		cmd, err := parseCmd(tt.text)
		if err != tt.err {
			t.Errorf("parseCmd(%q) error = \"%v\", want \"%v\"", tt.text, err, tt.err)
		}
		if cmd != tt.cmd {
			t.Errorf("parseCmd(%q) = %v, want %v", tt.text, cmd, tt.cmd)
		}
	}
}
