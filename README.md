# cgpcli

`cgpcli` is a high-performance Go parser and marshaler for **CommuniGate Pro** CLI. It provides a reliable bridge between CGP's unique syntax and native Go types with a focus on efficiency and correct type mapping.

## Data Representation: `map[string]any`

The library uses `map[string]any` as the primary container for CGP data structures.

* **Dynamic Schema:** CommuniGate Pro responses (e.g., `GetAccountSettings`) often contain hundreds of keys that vary between versions. Maps provide the flexibility to handle this without maintaining fragile, exhaustive structs.
* **Smart Typing:** The library performs type conversion following the CGP documentation:
    * **Standard Types:** Automatic conversion for Integers, Booleans, Arrays and Dictionaries.
    * **Specialized Types:** Values like `Created` dates, `Sizes` (10K, 5M), `Durations` (10m, never), and `IP addresses` are automatically converted to native Go types (`time.Time`, `int64`, `time.Duration`, `net.IP`).
    * **Legacy & Exceptions:** Correct handling of legacy formats and special CGP values through an internal registry.

## Application Architecture

For production applications, we recommend writing **thin wrappers** over the library's maps. Since a typical application only needs a fraction of the available keys, this approach ensures strict typing and clean code without the overhead of maintaining massive struct definitions.

> **Note:** Reflect-based struct mapping is planned for future updates to support interactive tools where development speed is prioritized over peak performance.

## Debugging

The library provides two levels of debugging to simplify development and troubleshooting:

* **Type Inspection:** The `PrettyPrint` function prints the structure of returned data to `stdout`. It highlights the Go types assigned to each value, clearly showing how the parser interpreted the CGP data.
* **Protocol Logging:** By calling `cli.SetDebug(true)`, you can enable raw protocol printing. This allows you to see the exact strings being exchanged with the CommuniGate Pro CLI, which is essential for diagnosing low-level communication issues.

## Performance & Design

* **Efficient Encoding:** The marshaling engine avoids unnecessary heap allocations and uses a specialized `writeInt64` implementation for maximum throughput.
* **Predictable Memory Usage:** Designed with careful resource management, avoiding unwarranted use of reflection and the `fmt` package in hot paths.

## License

This project is licensed under the BSD 3-Clause License.
