# go-pdebug

[![Build Status](https://travis-ci.org/lestrrat/go-pdebug.svg?branch=master)](https://travis-ci.org/lestrrat/go-pdebug)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-pdebug?status.svg)](https://godoc.org/github.com/lestrrat/go-pdebug)

Utilities for my print debugging fun. YMMV

# Synopsis

![optimized](https://pbs.twimg.com/media/CbiqhzLUUAIN_7o.png)

# Description

Building with `pdebug` declares a constant, `pdebug.Enabled` which you
can use to easily compile in/out depending on the presence of a build tag.

```go
func Foo() {
  // will only be available if you compile with `-tags debug`
  if pdebug.Enabled {
    pdebug.Printf("Starting Foo()!
  }
}
```

Note that using `github.com/lestrrat/go-pdebug` and `-tags debug` only
compiles in the code. In order to actually show the debug trace, you need
to specify an environment variable:

```shell
# For example, to show debug code during testing:
PDEBUG_ENABLE=1 go test -tags debug
```

If you want to forcefully show the trace (which is handy when you're
debugging/testing), you can use the `debug0` tag instead:

```shell
go test -tags debug0
```
