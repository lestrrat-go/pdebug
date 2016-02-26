//+build !debug,!debug0

package pdebug

// Enabled is true if `-tags debug` is used during compilation.
// Use this to basically "ifdef-out" debug blocks.
const Enabled = false

type guard struct{}

// IRelease undoes the indentation formatting. See IPrintf.
// IRelease is a no op unless you compile with the `debug` tag.
func (g guard) IRelease(f string, args ...interface{}) {}

// IPrintf prints out a message, and for subsequent calls,
// IPrintf is no op unless you comple with the `debug` tag.
func IPrintf(f string, args ...interface{}) guard { return guard{} }

// Printf prints to standard out, just like a normal fmt.Printf,
// but respects the indentation level set by IPrintf/IRelease.
// Printf is no op unless you compile with the `debug` tag.
func Printf(f string, args ...interface{}) {}

// Dump dumps the objects using go-spew.
// Dump is a no op unless you compile with the `debug` tag.
func Dump(v ...interface{}) {}
