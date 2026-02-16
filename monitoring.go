// # Monitoring and Server Control
//
// The Monitoring section provides access to the CommuniGate Pro internal
// statistics subsystem (SNMP-compatible). It enables real-time tracking of
// server performance, resource utilization, and protocol-specific metrics.
//
// Key capabilities include:
//   - Statistics Retrieval: Iterating through the hierarchical OID tree using
//     [Cli.GetNextStatName] and [Cli.ListAllStats] to discover available metrics.
//   - Value Inspection: Fetching current values for specific statistics elements
//     (counters, gauges, or strings) via [Cli.GetStatElement].
//   - Custom Metrics: Updating user-defined statistics elements using [Cli.SetStatElement]
//     with "SET" or "INC" (increment) operations.
//   - Signal Inspection: Retrieving detailed state information for active
//     Signal Dialogs by their numeric identifiers.
//   - Server Lifecycle: Gracefully initiating server termination via the [Cli.Shutdown] method.
//
// This section is vital for building health-check exporters (like Prometheus exporters),
// automated scaling triggers, and administrative dashboards.
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
//   - dialogID: The numeric Dialog ID.
//
// Returns:
//   - map[string]any: A dictionary with the Dialog status data.
//   - error: An error if the command fails.
func (cli *Cli) GetDialogInfo(dialogID int64) (map[string]any, error) {
	return cli.getMapAny("GETDIALOGINFO", dialogID)
}

// GetNextStatName is used to enumerate available Server statistics (SNMP) elements.
//
// Parameters:
//   - objectID: An empty string to get the first element, or the ID of the
//     previously found element.
//
// Returns:
//   - string: The ObjectID of the next statistics element.
//   - error: Returns io.EOF if the current element is the last one available.
func (cli *Cli) GetNextStatName(objectID string) (s string, err error) {
	s, err = cli.getString("GETNEXTSTATNAME", objectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", io.EOF
		}
		return "", err
	}
	return
}

// GetStatElement retrieves the current value of a Server statistics (SNMP) element.
//
// Parameters:
//   - objectID: The object ID of the Server statistics element.
//
// Returns:
//   - any: A number, string, or other object representing the element value.
//   - error: An error if the ObjectID is not found.
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
//   - []string: A slice containing all discovered statistic names.
//   - error: An error if any part of the iteration fails.
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
//   - objectID: The object ID of the Server statistics element.
//   - op: The operation to perform: "INC" (increment) or "SET" (assignment).
//   - val: The numeric value to be used in the operation.
//
// Returns:
//   - error: An error if the element is not "Custom" or the command fails.
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
// Returns:
//   - error: An error if the command fails.
func (cli *Cli) Shutdown() error {
	return cli.QueryNV("SHUTDOWN")
}
