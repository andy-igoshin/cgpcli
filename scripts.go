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

// RunScript executes a synchronous script (CG/PL) on the server.
// Unlike PBX tasks, this command blocks until the script completes
// and returns the result object directly.
//
// Parameters:
//   - account: The name of the Account. The script is started on this Account behalf.
//     The name can include the Domain name. If the Domain name is not specified, the current user Domain is used.
//   - prog: The name of the script file (the .scgp file) to start.
//   - entry: (Optional) The script entry point. If empty, "main" is used.
//   - param: (Optional) A map of parameters passed to the script.
//     Accessible via Vars().startParameter in CG/PL.
//
// Returns:
//   - any: The object returned by the script (string, map, array, etc.).
//   - error: An error if the command fails or the script execution encounters an error.
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
