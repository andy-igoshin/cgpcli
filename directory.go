// # Directory Administration
//
// The Directory section provides functionality for managing CommuniGate Pro Directory
// Units and their access control lists (ACLs). It enables integration with
// LDAP-compatible data stores and manages the hierarchical structure of the
// system-wide directory.
//
// Key capabilities include:
//   - Unit Lifecycle: Creating, relocating, and deleting local or remote
//     Directory Units with support for shared and mount-point configurations.
//   - Access Control: Granular management of directory access rights via the
//     [DirAccessRight] structure, including record-level and attribute-level permissions.
//   - Search and Discovery: Listing and retrieving configuration details for
//     individual directory units.
//
// Use these methods to synchronize account data, manage shared contacts,
// and configure organizational units within the CommuniGate Pro environment.
package cgpcli

import (
	"fmt"
)

// DirAccessRight defines a single rule for Directory Access Control.
type DirAccessRight struct {
	Name         string   // Unique name of the access rule.
	Target       string   // Target attribute or "*" for all attributes.
	BindDN       string   // Subject Distinguished Name (DN) this rule applies to.
	Type         string   // Access type: "allow", "prohibit", or "grant".
	RecordRights []string // Record-level permissions (e.g., "create", "delete").
	SearchAttrs  []string // List of attributes allowed for searching.
	ReadAttrs    []string // List of attributes allowed for reading.
	ModifyAttrs  []string // List of attributes allowed for modification.
}

// CreateDirectoryUnit creates a new Directory Unit.
//
// Parameters:
//   - unit: This parameter specifies the new Unit name.
//   - mountPoint: This parameter specifies the new Unit mount point (mount DN).
//   - shared: If true, a cluster-wide Directory Unit is created.
//   - remote: If true, a Remote (LDAP-based) Directory Unit is created,
//     otherwise a Local (File-based) Directory Unit is created.
//
// Returns:
//   - error: An error if the unit name or mount point is empty, or if
//     the server fails to create the unit.
func (cli *Cli) CreateDirectoryUnit(unit, mountPoint string, shared, remote bool) error {
	if unit == "" || mountPoint == "" {
		return fmt.Errorf("unit name and mount point are required")
	}

	args := make([]any, 0, 4)
	args = append(args, unit)
	if shared {
		args = append(args, "SHARED")
	}
	if remote {
		args = append(args, "REMOTE")
	}
	args = append(args, mountPoint)
	return cli.QueryNV("CREATEDIRECTORYUNIT", args...)
}

// DeleteDirectoryUnit removes an existing Directory Unit.
//
// Parameters:
//   - unit: This parameter specifies the Directory Unit name.
//   - shared: If true, this is a cluster-wide Directory Unit name.
//
// Returns:
//   - error: An error if the unit name is empty or if the server
//     fails to execute the command.
func (cli *Cli) DeleteDirectoryUnit(unit string, shared bool) error {
	if unit == "" {
		return fmt.Errorf("unit name is required")
	}

	const cmd = "DELETEDIRECTORYUNIT"
	if shared {
		return cli.QueryNV(cmd, unit, "SHARED")
	} else {
		return cli.QueryNV(cmd, unit)
	}
}

// GetDirectoryAccessRights retrieves the Directory Access Rights.
//
// Parameters:
//   - shared: If true, the cluster-wide Access Rights are retrieved.
//
// Returns:
//   - [][DirAccessRight]: A slice of Access Rights rules.
//   - error: An error if the server fails to retrieve or parse the data.
func (cli *Cli) GetDirectoryAccessRights(shared bool) ([]DirAccessRight, error) {
	const cmd = "GETDIRECTORYACCESSRIGHTS"
	var res []any
	var err error

	if shared {
		res, err = cli.getSliceAny(cmd, "SHARED")
	} else {
		res, err = cli.getSliceAny(cmd)
	}
	if err != nil {
		return nil, err
	}

	rights := make([]DirAccessRight, 0, len(res))
	for i, item := range res {
		fields, ok := item.([]any)
		if !ok || len(fields) < 8 {
			return nil, fmt.Errorf("%s: rule [%d] has invalid structure", cmd, i)
		}

		var r DirAccessRight
		if r.Name, err = toString(cmd, fields[0]); err != nil {
			return nil, err
		}
		if r.Target, err = toString(cmd, fields[1]); err != nil {
			return nil, err
		}
		if r.BindDN, err = toString(cmd, fields[2]); err != nil {
			return nil, err
		}
		if r.Type, err = toString(cmd, fields[3]); err != nil {
			return nil, err
		}
		if r.RecordRights, err = toStringSlice(cmd, fields[4]); err != nil {
			return nil, err
		}
		if r.SearchAttrs, err = toStringSlice(cmd, fields[5]); err != nil {
			return nil, err
		}
		if r.ReadAttrs, err = toStringSlice(cmd, fields[6]); err != nil {
			return nil, err
		}
		if r.ModifyAttrs, err = toStringSlice(cmd, fields[7]); err != nil {
			return nil, err
		}

		rights = append(rights, r)
	}
	return rights, nil
}

// GetDirectoryUnit retrieves the settings for a specific Directory Unit.
//
// Parameters:
//   - unit: This parameter specifies the Directory Unit name.
//   - shared: If true, this is a cluster-wide Directory Unit name.
//
// Returns:
//   - map[string]any: A dictionary with the Directory Unit settings.
//   - error: An error if the unit name is empty or the server fails
//     to retrieve the data.
func (cli *Cli) GetDirectoryUnit(unit string, shared bool) (map[string]any, error) {
	if unit == "" {
		return nil, fmt.Errorf("unit name is required")
	}

	const cmd = "GETDIRECTORYUNIT"
	if shared {
		return cli.getMapAny(cmd, unit, "SHARED")
	} else {
		return cli.getMapAny(cmd, unit)
	}
}

// ListDirectoryUnits retrieves the list of all created Directory units.
//
// Parameters:
//   - shared: If true, the cluster-wide Units are listed.
//
// Returns:
//   - map[string]string: A dictionary where keys are Directory Unit mount points
//     and values are Directory Unit names.
//   - error: An error if the server fails to execute the command.
func (cli *Cli) ListDirectoryUnits(shared bool) (map[string]string, error) {
	const cmd = "LISTDIRECTORYUNITS"
	var res map[string]any
	var err error

	if shared {
		res, err = cli.getMapAny(cmd, "SHARED")
	} else {
		res, err = cli.getMapAny(cmd)
	}
	if err != nil {
		return nil, err
	}

	units := make(map[string]string, len(res))
	for mnt, nameAny := range res {
		name, _ := toString(cmd, nameAny)
		units[name] = mnt
	}
	return units, nil
}

// RelocateDirectoryUnit re-mounts an existing Directory Unit on a different mount point.
//
// Parameters:
//   - unit: This parameter specifies the Directory Unit name.
//   - mountPoint: This parameter specifies the new mount point (mount DN).
//   - shared: If true, this is a cluster-wide Directory Unit name.
//
// Returns:
//   - error: An error if the unit name or mount point is empty, or if
//     the server fails to relocate the unit.
func (cli *Cli) RelocateDirectoryUnit(unit, mountPoint string, shared bool) error {
	if unit == "" || mountPoint == "" {
		return fmt.Errorf("unit name and new mount point are required")
	}

	const cmd = "RELOCATEDIRECTORYUNIT"
	if shared {
		return cli.QueryNV(cmd, unit, "SHARED", mountPoint)
	} else {
		return cli.QueryNV(cmd, unit, mountPoint)
	}
}

// SetDirectoryAccessRights sets the Directory Access Rights.
//
// Parameters:
//   - shared: If true, the cluster-wide Access Rights are set.
//   - rights: A slice of [DirAccessRight] structures specifying the new ACLs.
//
// Returns:
//   - error: An error if the rights array is nil or if the server
//     fails to apply the new settings.
func (cli *Cli) SetDirectoryAccessRights(shared bool, rights []DirAccessRight) error {
	if rights == nil {
		return fmt.Errorf("rights array is required")
	}

	const cmd = "SETDIRECTORYACCESSRIGHTS"
	ar := make([]any, len(rights))
	for i, r := range rights {
		ar[i] = []any{
			r.Name, r.Target, r.BindDN, r.Type,
			r.RecordRights, r.SearchAttrs, r.ReadAttrs, r.ModifyAttrs,
		}
	}

	if shared {
		return cli.QueryNV(cmd, "SHARED", ar)
	} else {
		return cli.QueryNV(cmd, ar)
	}
}

// SetDirectoryUnit changes the Directory Unit settings.
//
// Parameters:
//   - unit: This parameter specifies the Directory Unit name.
//   - shared: If true, this is a cluster-wide Directory Unit name.
//   - settings: A dictionary specifying the new Directory Unit settings.
//
// Returns:
//   - error: An error if the unit name or settings dictionary is empty,
//     or if the server fails to update the unit.
func (cli *Cli) SetDirectoryUnit(unit string, shared bool, settings map[string]any) error {
	if unit == "" || settings == nil {
		return fmt.Errorf("unit name and settings dictionary are required")
	}

	const cmd = "SETDIRECTORYUNIT"
	if shared {
		return cli.QueryNV(cmd, unit, "SHARED", settings)
	} else {
		return cli.QueryNV(cmd, unit, settings)
	}
}
