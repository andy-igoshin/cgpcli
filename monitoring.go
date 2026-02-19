// # Monitoring and Server Control
//
// The Monitoring section provides access to the CommuniGate Pro internal
// statistics subsystem (SNMP-compatible). It enables real-time tracking of
// server performance, resource utilization, and protocol-specific metrics.
//
// Key capabilities include:
//   - Statistics Retrieval: iterating through the hierarchical OID tree using
//     [Cli.GetNextStatName] and [Cli.ListAllStats] to discover available metrics.
//   - Value Inspection: fetching current values for specific statistics elements
//     (counters, gauges, or strings) via [Cli.GetStatElement].
//   - Custom Metrics: updating user-defined statistics elements using [Cli.SetStatElement]
//     with "SET" or "INC" (increment) operations.
//   - Signal Inspection: retrieving detailed state information for active
//     Signal Dialogs by their numeric identifiers.
//   - Server Lifecycle: gracefully initiating server termination via the [Cli.Shutdown] method.
package cgpcli

import (
	"fmt"
	"io"
	"strings"
)

// StatsCapacity определяет начальную емкость слайса для списка статистики.
// Значение 2048 покрывает стандартный набор OID большинства инсталляций без переаллокаций.
const statsCapacity = 2048

// GetDialogInfo retrieves the information about a Signal Dialog object.
//
// Parameters:
//   - dialogID: the numeric Dialog ID.
//
// This method executes the GETDIALOGINFO CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Dialog status data.
//   - error: an error if the command fails.
func (cli *Cli) GetDialogInfo(dialogID int64) (map[string]any, error) {
	return cli.getMapAny("GETDIALOGINFO", dialogID)
}

// GetNextStatName is used to enumerate available Server statistics (SNMP) elements.
//
// Parameters:
//   - objectID: an empty string to get the first element, or the ID of the
//     previously found element.
//
// This method executes the GETNEXTSTATNAME CLI command.
//
// Returns:
//   - string: the ObjectID of the next statistics element.
//   - error: returns io.EOF if the current element is the last one available.
func (cli *Cli) GetNextStatName(objectID string) (string, error) {
	s, err := cli.getString("GETNEXTSTATNAME", objectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", io.EOF
		}
		return "", err
	}
	return s, err
}

// GetStatElement retrieves the current value of a Server statistics (SNMP) element.
//
// Parameters:
//   - objectID: the object ID of the Server statistics element.
//
// This method executes the GETSTATELEMENT CLI command.
//
// Returns:
//   - any: a number, string, or other object representing the element value.
//   - error: an error if the command fails.
func (cli *Cli) GetStatElement(objectID string) (any, error) {
	if objectID == "" {
		return nil, fmt.Errorf("objectID is required")
	}
	return cli.Query("GETSTATELEMENT", objectID)
}

// ListAllStats is a helper method that iterates through the entire statistics tree
// and returns a slice of all available OIDs.
//
// Returns:
//   - []string: a slice containing all discovered statistic names.
//   - error: an error if the command fails.
func (cli *Cli) ListAllStats() ([]string, error) {
	stats := make([]string, 0, statsCapacity)
	var current string

	for {
		next, err := cli.GetNextStatName(current)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		stats = append(stats, next)
		current = next
	}
	return stats, nil
}

// SetStatElement updates the current value of a "Custom" Server statistics element.
//
// Parameters:
//   - objectID: the object ID of the Server statistics element.
//   - op: the operation to perform: "INC" (increment) or "SET" (assignment).
//   - val: the numeric value to be used in the operation.
//
// This method executes the SETSTATELEMENT CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetStatElement(objectID, op string, val int64) error {
	if objectID == "" {
		return fmt.Errorf("objectID is required")
	}
	if op != "INC" && op != "SET" {
		return fmt.Errorf("operation must be INC or SET")
	}
	return cli.QueryNV("SETSTATELEMENT", objectID, op, val)
}

// Shutdown initiates the stop sequence for the CommuniGate Pro Server.
//
// This method executes the SHUTDOWN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) Shutdown() error {
	return cli.QueryNV("SHUTDOWN")
}
