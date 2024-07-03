# xretry

The `xretry` package provides a mechanism to retry operations in your Go applications. It allows you to define a retry policy and use it to automatically retry a function if it fails.

## Features

- Define a custom retry policy with immediate retries and retries with backoff.
- Retry any function that returns an error.
- Backoff factor for exponential backoff in retries.

## Usage

Here's an example of how to use the `xretry` package:

```go
package main

import (
  "fmt"
  "github.com/hardiksachan/x/xretry"
  "time"
)

func main() {
  // Define a retry policy
  policy := xretry.NewRetryPolicy(
    xretry.WithImmediateRetries(3),
    xretry.WithRetriesWithBackoff(3, 1*time.Second, 2.0),
  )

  // Create a new retrier with the policy
  retrier := xretry.NewRetrier(policy)

  // Define a function that may fail
  f := func() error {
    // Your code here
  }

  // Use the retrier to retry the function if it fails
  err := retrier.Retry(f)
  if err != nil {
    fmt.Println("Operation failed after retries:", err)
  }
}
```

In this example, the function `f` will be retried immediately 3 times if it fails. If it still fails after these retries, it will be retried 3 more times with a delay that doubles after each retry, starting from 1 second.

## Installation

To use the `xretry` package in your Go project, you need to download it using the `go get` command:

```bash
go get github.com/hardiksachan/x
```

For more details about the `xretry` package, please refer to the source code in `xretry/retrier.go`.
