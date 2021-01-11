// +build !debug0,!debug

package pdebug_test

import (
	"testing"

	"github.com/lestrrat-go/pdebug/v3"
	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	// If we fail to provide this API, this test should fail to compile
	_ = pdebug.Marker
	_ = pdebug.FuncMarker
	_ = pdebug.Printf
}

func TestDisabled(t *testing.T) {
	if !assert.False(t, pdebug.Enabled) {
		return
	}
}
