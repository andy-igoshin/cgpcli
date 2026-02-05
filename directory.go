// # Section: Directory Administration
package cgpcli

import (
	"fmt"
)

type DirAccessRight struct {
	Name         string   // GUI: Name
	Target       string   // GUI: Target (attribute or *)
	BindDN       string   // GUI: Bind DN (subject)
	Type         string   // GUI: Type (allow, prohibit, grant)
	RecordRights []string // GUI: Record-Level Rights (create, delete)
	SearchAttrs  []string // GUI: Attribute Searching
	ReadAttrs    []string // GUI: Attribute Reading
	ModifyAttrs  []string // GUI: Attribute Modifying
}

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

func (cli *Cli) DeleteDirectoryUnit(unit string, shared bool) (err error) {
	if unit == "" {
		return fmt.Errorf("unit name is required")
	}

	const cmd = "DELETEDIRECTORYUNIT"
	if shared {
		err = cli.QueryNV(cmd, unit, "SHARED")
	} else {
		err = cli.QueryNV(cmd, unit)
	}
	return
}

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

func (cli *Cli) GetDirectoryUnit(unit string, shared bool) (res map[string]any, err error) {
	if unit == "" {
		return nil, fmt.Errorf("unit name is required")
	}

	const cmd = "GETDIRECTORYUNIT"
	if shared {
		res, err = cli.getMapAny(cmd, unit, "SHARED")
	} else {
		res, err = cli.getMapAny(cmd, unit)
	}
	return
}

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

func (cli *Cli) RelocateDirectoryUnit(unit, mountPoint string, shared bool) (err error) {
	if unit == "" || mountPoint == "" {
		return fmt.Errorf("unit name and new mount point are required")
	}

	const cmd = "RELOCATEDIRECTORYUNIT"
	if shared {
		err = cli.QueryNV(cmd, unit, "SHARED", mountPoint)
	} else {
		err = cli.QueryNV(cmd, unit, mountPoint)
	}
	return
}

func (cli *Cli) SetDirectoryAccessRights(shared bool, rights []DirAccessRight) (err error) {
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
		err = cli.QueryNV(cmd, "SHARED", ar)
	} else {
		err = cli.QueryNV(cmd, ar)
	}
	return
}

func (cli *Cli) SetDirectoryUnit(unit string, shared bool, settings map[string]any) (err error) {
	if unit == "" || settings == nil {
		return fmt.Errorf("unit name and settings dictionary are required")
	}

	const cmd = "SETDIRECTORYUNIT"
	if shared {
		err = cli.QueryNV(cmd, unit, "SHARED", settings)
	} else {
		err = cli.QueryNV(cmd, unit, settings)
	}
	return
}
