// # Section: Synchronous Scripts
package cgpcli

import (
	"fmt"
)

func (cli *Cli) RunScript(account, prog, entry string, param map[string]any) (any, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}
	if prog == "" {
		return nil, fmt.Errorf("program name is required")
	}

	args := make([]any, 0, 8)
	args = append(args, "RUNSCRIPT", account, "PROGRAM", prog)
	if entry != "" {
		args = append(args, "ENTRY", entry)
	}
	if param != nil {
		args = append(args, "PARAM", param)
	}

	return cli.Query(args...)
}
