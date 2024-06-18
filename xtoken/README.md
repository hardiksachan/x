# xtoken

xtoken is a Go package that provides token generation and handling functionalities. It is part of the larger X project, which provides a collection of libraries for various functionalities.

## Features

- Token Maker Interface: xtoken introduces a Maker interface that defines the methods for creating and verifying tokens.

## Usage

Here's an example of how to use xtoken:

```go
package main

import (
  "fmt"
  "time"
  "github.com/hardiksachan/x/xtoken"
)

func main() {
  var maker xtoken.Maker
  // Initialize the maker
  userID := "user1"
  email := "user1@example.com"
  duration := time.Minute

  token, payload, err := maker.CreateToken(userID, email, duration)
  if err != nil {
    fmt.Println("Error creating token:", err)
    return
  }

  fmt.Println("Created token:", token)
  fmt.Println("Token payload:", payload)

  verifiedPayload, err := maker.VerifyToken(token)
  if err != nil {
    fmt.Println("Error verifying token:", err)
    return
  }

  fmt.Println("Verified payload:", verifiedPayload)
}
```

For more details about each function and type, please refer to the source code in `xtoken/maker.go`.

---

Happy token handling!
