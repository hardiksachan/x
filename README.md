# X

This project provides a collection of libraries for various functionalities. Here's a brief overview of each library:

- [xmessage](xmessage/README.md) - provides asynchronous message handling functionalities, including idempotent inbox and transactional outbox.

- [xretry](xretry/README.md) - provides functionalities to retry operations.

- [xerrors](xerrors/README.md) - provides enhanced error handling capabilities.

- [xlog](xlog/README.md) - provides logging functionalities.

- [xtoken](xtoken/README.md) - provides token generation and handling functionalities.

- [xhash](xhash/README.md) - provides hashing functionalities, including bcrypt hashing.

- [xtest](xtest/README.md) - provides utilities for testing.

For more details about each library, please refer to their respective `README.md` files.

## Getting Started

These instructions will guide you on how to use this library in your Go projects.

### Installation

To use this library in your Go project, you need to download it using the `go get` command:

```bash
go get github.com/hardiksachan/x
```

### Usage

After downloading and importing the library, you can use its functionalities in your Go code. Here's an example:

```go
package main

import (
    "fmt"
    "github.com/hardiksachan/x/xerrors"
)

func main() {
    err := xerrors.New("this is an error")
    fmt.Println(err)
}
```

## Author

**Hardik Sachan**

- [Github](https://github.com/hardiksachan)
- [LinkedIn](https://www.linkedin.com/in/hardik-sachan/)
- [X](https://x.com/hardik__sachan)

## Support

Give a ⭐️ if you like this project!
