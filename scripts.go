// # Synchronous Script Execution
//
// The Scripts section provides methods for executing CommuniGate Pro
// Programming Language (CG/PL) scripts in a synchronous manner. Unlike
// standard PBX tasks that run asynchronously, these methods block until
// the script completes, returning the execution result directly to the caller.
//
// Key capabilities include:
//   - Immediate Execution: Using [Cli.RunScript] to trigger a specific program
//     and entry point within the context of an account.
//   - Data Exchange: Passing complex parameters via a map and receiving
//     structured results (any type supported by CommuniGate Pro CLI) back into the Go application.
//   - Contextual Execution: Scripts run within the environment of a specified
//     account, respecting its settings and storage access.
//
// This section is ideal for performing complex calculations, data
// transformations, or administrative tasks that require immediate feedback
// and logic implemented in CG/PL.
package cgpcli

import (
	"fmt"
)

func (cli *Cli) RunScript(account, prog, entry string, param map[string]any) (any, error) {
	if account == "" || prog == "" {
		return nil, fmt.Errorf("account name and prog name are required")
	}

	args := make([]any, 0, 7)
	args = append(args, account, "PROGRAM", prog)
	if entry != "" {
		args = append(args, "ENTRY", entry)
	}
	if param != nil {
		args = append(args, "PARAM", param)
	}
	return cli.Query("RUNSCRIPT", args...)
}
