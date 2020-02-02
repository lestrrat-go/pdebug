// +build debug OR debug0

package pdebug

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

const Enabled = true

type markerKey struct{}
type markerCtx struct {
	clock      interface{ Now() time.Time }
	indent     int
	out        io.Writer
	prefix     []byte
	timestamps bool
}

var defaultClock = ClockFunc(time.Now)

func defaultMarkerCtx() *markerCtx {
	return &markerCtx{
		clock:      defaultClock,
		out:        os.Stderr,
		prefix:     []byte("|DEBUG| "),
		timestamps: true,
	}
}

func getMarkerCtx(ctx context.Context) (context.Context, *markerCtx) {
	var mctx *markerCtx
	v := ctx.Value(markerKey{})
	if v == nil {
		mctx = defaultMarkerCtx()
		ctx = context.WithValue(ctx, markerKey{}, mctx)
	} else {
		mctx = v.(*markerCtx)
	}

	return ctx, mctx
}

func WithTimestamp(ctx context.Context, b bool) context.Context {
	if !Trace {
		return ctx
	}

	xctx, mctx := getMarkerCtx(ctx)
	mctx.timestamps = b
	return xctx
}

func WithClock(ctx context.Context, clock interface{ Now() time.Time }) context.Context {
	if !Trace {
		return ctx
	}

	xctx, mctx := getMarkerCtx(ctx)
	mctx.clock = clock
	return xctx
}

func WithOutput(ctx context.Context, out io.Writer) context.Context {
	if !Trace {
		return ctx
	}

	xctx, mctx := getMarkerCtx(ctx)
	mctx.out = out
	return xctx
}

func WithPrefix(ctx context.Context, prefix string) context.Context {
	if !Trace {
		return ctx
	}

	xctx, mctx := getMarkerCtx(ctx)
	mctx.prefix = []byte(prefix)
	return xctx
}

type MarkerGuard struct {
	ctx       context.Context
	errptr    *error
	indent    int
	msgFormat string
	msgArgs   []interface{}
	out       io.Writer
	prefix    []byte
	start     time.Time
}

func formatMarkerMessage(buf *[]byte, format string, args []interface{}, prefix, postfix []byte, clock Clock, indent int) {
	// foo\nbar\n should be written as preamble foo\npreamble bar\n
	var scratch bytes.Buffer
	fmt.Fprintf(&scratch, format, args...)

	scanner := bufio.NewScanner(&scratch)
	for scanner.Scan() {
		appendPreamble(buf, prefix, clock, indent)
		line := scanner.Bytes()
		*buf = append(*buf, line...)
		*buf = append(*buf, postfix...)
		*buf = append(*buf, '\n')
	}
}

var markerGuardPool = sync.Pool{
	New: allocMarkerGuard,
}

func allocMarkerGuard() interface{} {
	return &MarkerGuard{}
}

func getMarkerGuard() *MarkerGuard {
	return markerGuardPool.Get().(*MarkerGuard)
}

func releaseMarkerGuard(mg *MarkerGuard) {
	mg.ctx = nil
	mg.indent = 0
	mg.msgFormat = ""
	mg.msgArgs = nil
	mg.prefix = nil
	markerGuardPool.Put(mg)
}

// Marker creates a marker. A marker is basically something that is used
// to remember and mark the entry point and the exit point of a particular
// section of code.
func Marker(ctx context.Context, format string, args ...interface{}) *MarkerGuard {
	if !Trace {
		return nil
	}

	xctx, mctx := getMarkerCtx(ctx)
	mg := getMarkerGuard()
	mg.ctx = xctx
	mg.indent = mctx.indent
	mg.msgFormat = format
	mg.msgArgs = args
	mg.prefix = mctx.prefix

	if clock := mctx.clock; clock != nil {
		mg.start = clock.Now()
	}

	mctx.indent += 2
	var buf []byte

	var clock Clock
	if mctx.timestamps {
		clock = mctx.clock
	}

	formatMarkerMessage(&buf, "START "+mg.msgFormat, mg.msgArgs, mg.prefix, nil, clock, mg.indent)
	mctx.out.Write(buf)
	return mg
}

func (mg *MarkerGuard) BindError(err *error) *MarkerGuard {
	if !Trace {
		return nil
	}

	mg.errptr = err
	return mg
}

func appendPreamble(buf *[]byte, prefix []byte, clock Clock, indent int) {
	*buf = append(*buf, prefix...)
	if clock != nil {
		*buf = append(*buf,
			[]byte(strconv.FormatFloat(float64(clock.Now().UnixNano())/1000000.0, 'f', 5, 64))...,
		)
		*buf = append(*buf, ' ')
	}
	for i := 0; i < indent; i++ {
		*buf = append(*buf, ' ')
	}
}

// End finalizes the MarkerGuard. Subsequent calls to the same object are
// invalid, and may cause panics.
func (mg *MarkerGuard) End() {
	if !Trace {
		return
	}

	_, mctx := getMarkerCtx(mg.ctx)
	if mctx.indent < 2 {
		mctx.indent = 0
	} else {
		mctx.indent -= 2
	}

	var postfix []byte
	var clock Clock
	if mctx.timestamps {
		clock = mctx.clock
	}
	if clock != nil || mg.errptr != nil {
		postfix = append(postfix, '(')
		if clock != nil {
			postfix = append(postfix, []byte("elapsed=")...)
			postfix = append(postfix, []byte(clock.Now().Sub(mg.start).String())...)
		}

		if errptr := mg.errptr; errptr != nil && *errptr != nil {
			if clock != nil {
				postfix = append(postfix, ", "...)
			}
			postfix = append(postfix, "error="...)
			postfix = append(postfix, []byte((*errptr).Error())...)
		}
		postfix = append(postfix, ')')
	}

	var buf []byte
	formatMarkerMessage(&buf, "END   "+mg.msgFormat, mg.msgArgs, mg.prefix, postfix, clock, mctx.indent)

	mctx.out.Write(buf)

	releaseMarkerGuard(mg)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	if !Trace {
		return
	}

	_, mctx := getMarkerCtx(ctx)

	var buf []byte
	var clock Clock
	if mctx.timestamps {
		clock = mctx.clock
	}
	formatMarkerMessage(&buf, format, args, mctx.prefix, nil, clock, mctx.indent)
	mctx.out.Write(buf)
}
