# xhash

`xhash` is a Go package that provides hashing functionalities. It is part of the larger X project, which provides a collection of libraries for various functionalities.

## Features

- **Bcrypt Hashing**: `xhash` provides a simple interface for hashing passwords using bcrypt. This includes functions to hash a password and to compare a hashed password with a plaintext one.

## Usage

Here's an example of how to use `xhash`:

```go
package main

import (
    "fmt"
    "github.com/hardiksachan/x/xhash"
)

func main() {
    password := "mysecretpassword"
    hashedPassword, err := xhash.HashPassword(password)
    if err != nil {
        fmt.Println("Error hashing password:", err)
        return
    }

    fmt.Println("Hashed password:", hashedPassword)

    err = xhash.ComparePassword(hashedPassword, password)
    if err != nil {
        fmt.Println("Error comparing password:", err)
        return
    }

    fmt.Println("Password comparison successful!")
}
```
