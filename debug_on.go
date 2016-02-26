// +build debug OR debug0

package pdebug

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const Enabled = true

var prefix = ""
var prefixToken = "  "
var logger = log.New(os.Stdout, "|DEBUG| ", 0)

type guard func() time.Time

func emptyGuard() time.Time {
	return time.Time{}
}

func (g guard) IRelease(f string, args ...interface{}) {
	if !Trace {
		return
	}
	start := g()
	dur := time.Since(start)
	Printf("%s (%s)", fmt.Sprintf(f, args...), dur)
}

// IPrintf indents and then prints debug messages. Execute the callback
// to undo the indent
func IPrintf(f string, args ...interface{}) guard {
	if !Trace {
		return emptyGuard
	}
	Printf(f, args...)
	prefix = prefix + prefixToken
	start := time.Now()
	return guard(func() time.Time {
		prefix = prefix[len(prefixToken):]
		return start
	})
}

// Printf prints debug messages. Only available if compiled with "debug" tag
func Printf(f string, args ...interface{}) {
	if !Trace {
		return
	}
	logger.Printf("%s%s", prefix, fmt.Sprintf(f, args...))
}

func Dump(v ...interface{}) {
	if !Trace {
		return
	}
	spew.Dump(v...)
}
