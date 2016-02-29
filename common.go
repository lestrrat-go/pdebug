package pdebug

import (
	"io"
	"os"
	"time"
)

type pdctx struct {
	indentL int
	Prefix  string
	Writer  io.Writer
}

var emptyMarkerGuard = &markerg{}

type markerg struct {
	indentg guard
	ctx     *pdctx
	f       string
	args    []interface{}
	start   time.Time
	errptr  *error
}

var DefaultCtx = &pdctx{
	Prefix: "|DEBUG| ",
	Writer: os.Stdout,
}

type guard struct {
	cb func()
}

func (g *guard) End() {
	if cb := g.cb; cb != nil {
		cb()
	}
}
