// # Section: Synchronous Scripts
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
