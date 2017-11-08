package stdio

import (
	"bytes"
	"fmt"
	"testing"
)

type testcase struct {
	name string
	qval bool
	vval bool
	wout string
	werr string
}

const (
	msg = "message"
	yes = msg + "\n"
	no  = ""
)

func (tc *testcase) test(t *testing.T) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	s := New()
	s.Verbose(&tc.vval)
	s.Quiet(&tc.qval)
	s.Stdout(stdout)
	s.Stderr(stderr)

	switch tc.name {
	case "Mention":
		s.Mention(msg)
	case "Warn":
		s.Warn(msg)
	case "Babble":
		s.Babble(msg)
	case "Caution":
		s.Caution(msg)
	}

	if got := stdout.String(); got != tc.wout {
		t.Errorf("with (verbose=%v, quiet=%v), %s(%q) wrote %q to stdout; Wanted %q", tc.vval, tc.qval, tc.name, got, tc.wout)
	}

	if got := stderr.String(); got != tc.werr {
		t.Errorf("with (verbose=%v, quiet=%v), %s(%q) wrote %q to stderr; Wanted %q", tc.vval, tc.qval, tc.name, got, tc.werr)
	}
}

func TestStdio(t *testing.T) {
	tests := []*testcase{
		// method   quiet verbose wout werr
		{"Mention", false, false, yes, no},
		{"Mention", false, true, yes, no},
		{"Mention", true, false, no, no},
		{"Mention", true, true, no, no},

		{"Warn", false, false, no, yes},
		{"Warn", false, true, no, yes},
		{"Warn", true, false, no, no},
		{"Warn", true, true, no, no},

		{"Babble", false, false, no, no},
		{"Babble", false, true, yes, no},
		{"Babble", true, false, no, no},
		{"Babble", true, true, no, no},

		{"Caution", false, false, no, no},
		{"Caution", false, true, no, yes},
		{"Caution", true, false, no, no},
		{"Caution", true, true, no, no},
	}

	for _, tc := range tests {
		label := fmt.Sprintf("method=%s:quiet=%v:verbose=%v", tc.name, tc.qval, tc.vval)
		t.Run(label, tc.test)
	}
}
