// Package errstack provides an error type that adds a stack trace to an error.
package errstack

import (
	"fmt"
	"log"
	"runtime"
)

// The size of the stack trace buffer for reading.  Stack traces will not be
// longer than this.
const StackTraceSize = 4096

// Wraps an error and attaches the stack trace from where you call
// errstack.New(err)
type ErrorStackTrace struct {
	// The original error
	Err error
	// The stack trace at the time when errstack.New(err) was called
	StackTrace string
}

// Satisfies the error interface
func (self *ErrorStackTrace) Error() string {
	return self.String()
}

func (self *ErrorStackTrace) String() string {
	return fmt.Sprintf("%v\n%s", self.Err, self.StackTrace)
}

// Attaches a stack trace to an error.  If the error is already an
// ErrorStackTrace, it is a no-op and returns the error given as argument.
// This is so you can call New(err) on all errors as they are passed up
// the call stack, without losing the stack trace for where the error
// originally occurred.
func New(err error) *ErrorStackTrace {
	if err == nil {
		return nil
	}
	if stackError, ok := err.(*ErrorStackTrace); ok {
		return stackError
	}
	var stackBuf [StackTraceSize]byte
	n := runtime.Stack(stackBuf[:], false)
	// The stack will look like:
	//     goroutine 1 [running]:
	//     errstack.New()
	//         /tmpfs/gosandbox-xxx/stackerror.go:32 +0xe0
	//     ... remaining stack ...
	// cut out the 2nd and 3rd lines
	stack := cutLines(stackBuf[:n], 1, 2)
	return &ErrorStackTrace{
		Err:        err,
		StackTrace: string(stack),
	}
}

// Removes lines between startLine and endLine, inclusive. i.e.
// startLine = 2 and startLine = 4 will remove lines 2, 3, and 4. Lines are
// 0-indexed.
func cutLines(str []byte, startLine, endLine int) []byte {
	if startLine < 0 || endLine < 0 || endLine <= startLine {
		log.Panic("Invalid startLine or endLine")
	}
	start := -1
	end := -1
	currLine := 0
	if startLine == 0 {
		start = 0
	}
	for i, _ := range str {
		if str[i] == '\n' {
			currLine += 1
			if start == -1 {
				if currLine == startLine {
					start = i + 1
				}
			} else {
				if currLine-1 == endLine {
					end = i + 1
					break
				}
			}
		}
	}
	if end < 0 {
		end = len(str)
	}
	if start >= 0 && end > start {
		str = append(str[:start], str[end:]...)
	}
	return str
}
