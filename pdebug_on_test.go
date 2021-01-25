// +build debug0 or debug

package pdebug_test

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/lestrrat-go/pdebug/v3"
	"github.com/stretchr/testify/assert"
)

func TestMarker(t *testing.T) {
	fn := func(t *testing.T, wg *sync.WaitGroup) {
		t.Helper()
		if wg != nil {
			defer wg.Done()
		}
		var err error
		g1 := pdebug.FuncMarker().BindError(&err)
		defer g1.End()

		pdebug.Printf("Hello, World test 1")

		g2 := pdebug.Marker("Test")
		defer g2.End()

		pdebug.Printf("Hello, World test 2")
		err = errors.New("test 1 error")
	}

	var buf bytes.Buffer
	var now = time.Unix(0, 0)
	pdebug.Configure(
		pdebug.WithClock(pdebug.ClockFunc(func() time.Time { return now })),
		pdebug.WithWriter(&buf),
	)

	t.Run("Output format", func(t *testing.T) {
		fn(t, nil)
		t.Logf("%s", buf.String())
		if pdebug.Enabled && pdebug.Trace {
			const expected = `|DEBUG| 0.00000 START github.com/lestrrat-go/pdebug/v3_test.TestMarker.func1
|DEBUG| 0.00000   Hello, World test 1
|DEBUG| 0.00000   START Test
|DEBUG| 0.00000     Hello, World test 2
|DEBUG| 0.00000   END   Test(elapsed=0s)
|DEBUG| 0.00000 END   github.com/lestrrat-go/pdebug/v3_test.TestMarker.func1(elapsed=0s, error=test 1 error)
`
			if !assert.Equal(t, expected, buf.String()) {
				return
			}
		}
	})
	t.Run("Race condition", func(t *testing.T) {
		// Make sure race conditions don't exist by calling multiple goroutines
		const gcount = 10

		var wg sync.WaitGroup
		wg.Add(gcount)
		for i := 0; i < gcount; i++ {
			go fn(t, &wg)
		}
		wg.Wait()
	})

}
