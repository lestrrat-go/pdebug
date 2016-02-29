package pdebug

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setw(ctx *pdctx, w io.Writer) func() {
	oldw := ctx.Writer
	ctx.Writer = w
	return func() { ctx.Writer = oldw }
}

func TestPrintf(t *testing.T) {
	buf := &bytes.Buffer{}
	wg := setw(DefaultCtx, buf)
	defer wg()

	Printf("Hello, World!")

	if Enabled && Trace {
		if !assert.Equal(t, "|DEBUG| Hello, World!\n", buf.String(), "Simple Printf works") {
			return
		}
	} else {
		if !assert.Equal(t, "", buf.String(), "Simple Printf should be supressed") {
			return
		}
	}
}

func TestMarker(t *testing.T) {
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

	if Enabled && Trace {
		const expected = "|DEBUG| START f1\n|DEBUG| \tSTART f2\n|DEBUG| \t\tHello, World!\n|DEBUG| \tEND f2 ("
		if !assert.True(t, strings.HasPrefix(buf.String(), expected), "Markers should work") {
			t.Logf("Expected '%v'", expected)
			t.Logf("Actual   '%v'", buf.String())
			return
		}
	} else {
		if !assert.Equal(t, "", buf.String(), "Markers should work") {
			return
		}
	}
}
