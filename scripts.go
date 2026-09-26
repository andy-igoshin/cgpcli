// # Synchronous Script Execution
//
// The Scripts section provides methods for executing CommuniGate Pro
// Programming Language (CG/PL) scripts in a synchronous manner. Unlike
// standard PBX tasks that run asynchronously, these methods block until
// the script completes, returning the execution result directly to the caller.
//
// Key capabilities include:
//   - Immediate Execution: using [Cli.RunScript] to trigger a specific program
//     and entry point within the context of an account.
//   - Data Exchange: passing complex parameters via a map and receiving
//     structured results (any type supported by CommuniGate Pro CLI) back into the Go application.
//   - Contextual Execution: scripts run within the environment of a specified
//     account, respecting its settings and storage access.
//
// This section is ideal for performing complex calculations, data
// transformations, or administrative tasks that require immediate feedback
// and logic implemented in CG/PL.
package cgpcli

// RunScript executes a synchronous script (CG/PL) on the server.
// Unlike PBX tasks, this command blocks until the script completes
// and returns the result object directly.
//
// Parameters:
//   - account: the Account name. The name can include the Domain name.
//     The script is executed on behalf of this Account.
//   - prog: the name of the script file (the .scgp file) to start.
//   - entry: an optional script entry point. If empty, "main" is used.
//   - param: an optional map of parameters passed to the script.
//     Accessible via Vars().startParameter in CG/PL.
//
// This method executes the PROGRAM CLI command.
//
// Returns:
//   - any: the object returned by the script (string, map, array, etc.).
//   - error: an error if the command fails.
func (cli *Cli) RunScript(account, prog, entry string, param map[string]any) (any, error) {
	if account == "" {
		return nil, ErrAccountRequired
	}
	if prog == "" {
		return nil, ErrProgramRequired
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
