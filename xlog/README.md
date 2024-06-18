## xlog

`xlog` is a Go package that provides logging functionalities. It is part of the larger X project, which provides a collection of libraries for various functionalities.

## Features

- **Log Interface**: `xlog` introduces a `Log` interface that defines the methods for logging messages at different levels: Info, Debug, Warn, and Error.

- **Message Struct**: `xlog` provides a `Message` struct that carries information about the log message, including title, details, and additional data.

- **Current Logger**: `xlog` provides a `currentLogger` function to get the current logger instance.

- **Noop Logger**: `xlog` provides a `noopLogger` that implements the `Log` interface but does nothing. This can be useful in testing or when you want to disable logging.

- **Pretty Logger**: `xlog` provides a `prettyLogger` that implements the `Log` interface and logs messages in a human-readable format. This uses `zap` under the hood.

## Usage

Here's an example of how to use `xlog`:

```go
package main

import (
    "github.com/hardiksachan/x/xlog"
)

func main() {
    msg := xlog.Message{
        Title: "Test Log",
        Details: "This is a test log message",
        Data: map[string]string{
            "key": "value",
        },
    }

    xlog.Info(msg)
}
```

For more details about each function and type, please refer to the source code in `xlog/current.go`, `xlog/log.go`, `xlog/noop.go`, and `xlog/pretty.go`.

---

Happy logging!
