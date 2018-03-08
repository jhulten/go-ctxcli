# ctxcli: Extend golang context to CLI [![GoDoc](https://godoc.org/github.com/jhulten/go-ctxcli?status.png)](https://godoc.org/github.com/jhulten/go-ctxcli) [![Build Status](https://travis-ci.org/jhulten/go-ctxcli.svg?branch=master)](https://travis-ci.org/jhulten/go-ctxcli)

ctxcli is a Go library extending the standard context with handling for OS signals. 

## Features

- Create a child context that cancels when an OS signal like SIGINT is received.
- Helper functions to exit or panic when the provided context is cancelled.
- Consistent API with the context standard library.
 
## Example


```golang
package main

import (
    "context"
    "log"
    "time"
    "github.com/jhulten/go-ctxcli"
)

func main() {
    ctx := context.Background()
    ctx = ctxcli.WithInterrupt(ctx)

    for {
        log.Print(".")
        time.Sleep(500 * time.Millisecond)
        ctxcli.ExitIfCancelled(ctx)
    }
}
```

