# cgpcli

[![Go Reference](https://pkg.go.dev/badge/git.vsu.ru/ai/cgpcli.svg)](https://pkg.go.dev/git.vsu.ru/ai/cgpcli)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](LICENSE.md)

`cgpcli` is a Go client library for the **CommuniGate Pro** CLI. It provides a reliable bridge between CGP's unique syntax and native Go types with a focus on efficiency and correct type mapping.

## Documentation

Full API documentation and usage examples are available at [pkg.go.dev](https://pkg.go.dev/git.vsu.ru/ai/cgpcli).

## Installation

```bash
go get git.vsu.ru/ai/cgpcli
```

## Usage

```go
package main

import (
	"git.vsu.ru/ai/cgpcli"
)

func main() {
	cli, err := cgpcli.New("mail.example.com", "admin", "password", cgpcli.Plain, false)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	settings, err := cli.GetAccountSettings("user@example.com")
	if err != nil {
		panic(err)
	}

	cgpcli.PrettyPrint(settings)
}
```

## Data Representation: `map[string]any`

The library uses `map[string]any` as the primary container for CGP data structures.

- **Dynamic Schema:** CommuniGate Pro responses (e.g., `GetAccountSettings`) often contain hundreds of keys that vary between versions. Maps provide the flexibility to handle this without maintaining large, fragile struct definitions.
- **Smart Typing:** The library performs type conversion following the CGP documentation:
  - **Standard Types:** Automatic conversion for integers, booleans, arrays, and dictionaries.
  - **Specialized Types:** Values like `Created` dates, sizes (`10K`, `5M`), durations (`10m`, `never`), and IP addresses are automatically converted to native Go types (`time.Time`, `int64`, `time.Duration`, `net.IP`).
  - **Legacy & Exceptions:** Correct handling of legacy formats and special CGP values through an internal registry.

## Application Architecture

For production applications, it is recommended to write **thin wrappers** over the library's maps. Since a typical application only needs a fraction of the available keys, this approach ensures strict typing and clean code without the overhead of maintaining comprehensive struct definitions.

## Debugging

The library provides two levels of debugging to simplify development and troubleshooting:

- **Type Inspection:** The `PrettyPrint` function prints the structure of returned data to `stdout`. It highlights the Go types assigned to each value, clearly showing how the parser interpreted the CGP data.
- **Protocol Logging:** By calling `cli.SetDebug(true)`, you can enable raw protocol logging to `stderr`. This allows you to see the exact strings being exchanged with the CommuniGate Pro CLI, which is essential for diagnosing low-level communication issues.

## Performance and Design

The marshaling engine avoids unnecessary heap allocations. The library avoids use of reflection and the `fmt` package in hot paths.

## License

This project is licensed under the BSD 3-Clause License.

## Author

Andrey Igoshin <ai@vsu.ru>

### Links

* Repository: <https://git.vsu.ru/ai/rspamd-cgp>
* CommuniGate Pro Website: <https://communigatepro.ru>
* Command Line Interface/API: <https://doc.communigatepro.ru/russian/development/CLI.html>
