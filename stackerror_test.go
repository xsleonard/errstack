package errstack

import (
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	err := errors.New("bad")
	stackErr := New(err)
	lines := strings.Split(stackErr.StackTrace, "\n")
	if len(lines) != 8 {
		t.Fatalf("Expected 8 lines of stack trace. Have %d", len(lines))
	}
	if stackErr.Err != err {
		t.Fatal("StackError.Err has different value")
	}
	if !strings.HasPrefix(stackErr.Error(), "bad") {
		t.Fatalf("error string does begin with error")
	}
	if stackErr.Error() != stackErr.String() {
		t.Fatalf("Error() != String()")
	}
}
