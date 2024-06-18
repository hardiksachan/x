# xerrors

`xerrors` is a Go package that provides enhanced error handling capabilities. It is part of the larger X project, which provides a collection of libraries for various functionalities.

## Features

- Custom Error Type: `xerrors` introduces a custom `Error` type that carries information about operation, error code, message, and the underlying error.

- Error Codes: `xerrors` provides a set of predefined error codes like `Other`, `Internal`, `Invalid`, `NotFound`, `Exists`, and `Expired`.

- Error Messages: `xerrors` allows you to associate human-readable messages with your errors.

- gRPC Error Handling: `xerrors` provides functions to convert errors to gRPC errors.

---

Happy `err` - ing!
