// # Directory Administration
//
// The Directory section provides functionality for managing CommuniGate Pro Directory
// Units and their access control lists (ACLs). It enables integration with
// LDAP-compatible data stores and manages the hierarchical structure of the
// system-wide directory.
//
// Key capabilities include:
//   - Unit Lifecycle: creating, relocating, and deleting local or remote
//     Directory Units with support for shared and mount-point configurations.
//   - Access Control: granular management of directory access rights via the
//     [DirAccessRight] structure, including record-level and attribute-level permissions.
//   - Search and Discovery: listing and retrieving configuration details for
//     individual directory units.
package cgpcli

import (
	"fmt"
)

// DirAccessRight defines a single rule for Directory Access Control.
type DirAccessRight struct {
	Name         string   // unique name of the access rule.
	Target       string   // target attribute, or "*" for all attributes.
	BindDN       string   // subject distinguished name (DN) this rule applies to.
	Type         string   // access type: "allow", "prohibit", or "grant".
	RecordRights []string // record-level permissions (e.g., "create", "delete").
	SearchAttrs  []string // list of attributes allowed for searching.
	ReadAttrs    []string // list of attributes allowed for reading.
	ModifyAttrs  []string // list of attributes allowed for modification.
}

// CreateDirectoryUnit creates a new Directory Unit.
//
// Parameters:
//   - unit: the new Unit name.
//   - mountPoint: the new Unit mount point (mount DN).
//   - shared: if true, a Cluster-wide Directory Unit is created.
//   - remote: if true, a Remote (LDAP-based) Directory Unit is created,
//     otherwise a Local (File-based) Directory Unit is created.
//
// This method executes the CREATEDIRECTORYUNIT CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
//   - unit: the Directory Unit name.
//   - shared: if true, this is a Cluster-wide Directory Unit name.
//
// This method executes the DELETEDIRECTORYUNIT CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
//   - shared: if true, the Cluster-wide Access Rights are retrieved.
//
// This method executes the GETDIRECTORYACCESSRIGHTS CLI command.
//
// Returns:
//   - [][DirAccessRight]: a slice of Access Rights rules.
//   - error: an error if the command fails.
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
//   - unit: the Directory Unit name.
//   - shared: if true, this is a Cluster-wide Directory Unit name.
//
// This method executes the GETDIRECTORYUNIT CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Directory Unit settings.
//   - error: an error if the command fails.
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
//   - shared: if true, the Cluster-wide Units are listed.
//
// This method executes the LISTDIRECTORYUNITS CLI command.
//
// Returns:
//   - map[string]string: a dictionary where keys are Directory Unit mount points
//     and values are Directory Unit names.
//   - error: an error if the command fails.
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
//   - unit: the Directory Unit name.
//   - mountPoint: the new mount point (mount DN).
//   - shared: if true, this is a Cluster-wide Directory Unit name.
//
// This method executes the RELOCATEDIRECTORYUNIT CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
//   - shared: if true, the Cluster-wide Access Rights are set.
//   - rights: a slice of [DirAccessRight] structures specifying the new ACLs.
//
// This method executes the SETDIRECTORYACCESSRIGHTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
//   - unit: the Directory Unit name.
//   - shared: if true, this is a Cluster-wide Directory Unit name.
//   - settings: a dictionary of Directory Unit settings.
//
// This method executes the SETDIRECTORYUNIT CLI command.
//
// Returns:
//   - error: an error if the command fails.
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
