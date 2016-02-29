//+build !debug0,!debug

package pdebug

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setw(ctx *pdctx, w io.Writer) func() {
	oldw := ctx.Writer
	ctx.Writer = w
	return func() { ctx.Writer = oldw }
}

func TestDebug0Basic(t *testing.T) {
	buf := &bytes.Buffer{}
	wg := setw(DefaultCtx, buf)
	defer wg()

	Printf("Hello, World!")

	if !assert.Equal(t, "", buf.String(), "Simple Printf works") {
		return
	}
}

func TestDebug0Indent(t *testing.T) {
	buf := &bytes.Buffer{}
	wg := setw(DefaultCtx, buf)
	defer wg()

	f2 := func() (err error) {
		g := Marker("f2").BindError(&err)
		defer g.End()
		Printf("Hello, World!")
		return errors.New("dummy error")
	}

	f1 := func() {
		g := Marker("f1")
		defer g.End()
		f2()
	}

	f1()

	if !assert.Equal(t, "", buf.String(), "Printf with indent works") {
		return
	}
}