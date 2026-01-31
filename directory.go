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

func (cli *Cli) CreateDirectoryUnit(name, mountPoint string, shared, remote bool) error {
	if name == "" || mountPoint == "" {
		return fmt.Errorf("unit name and mount point are required")
	}

	args := make([]any, 0, 4)
	args = append(args, name)
	if shared {
		args = append(args, "SHARED")
	}
	if remote {
		args = append(args, "REMOTE")
	}
	args = append(args, mountPoint)
	return cli.QueryNV("CREATEDIRECTORYUNIT", args...)
}

func (cli *Cli) DeleteDirectoryUnit(name string, shared bool) (err error) {
	if name == "" {
		return fmt.Errorf("unit name is required")
	}
	const cmd = "DELETEDIRECTORYUNIT"
	if shared {
		err = cli.QueryNV(cmd, name, "SHARED")
	} else {
		err = cli.QueryNV(cmd, name)
	}
	return
}

func (cli *Cli) GetDirectoryAccessRights(shared bool) ([]DirAccessRight, error) {
	const cmd = "GETDIRECTORYACCESSRIGHTS"
	var res any
	var err error

	if shared {
		res, err = cli.Query(cmd, "SHARED")
	} else {
		res, err = cli.Query(cmd)
	}
	if err != nil {
		return nil, err
	}

	raw, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected []any, got %T", cmd, res)
	}

	rights := make([]DirAccessRight, 0, len(raw))
	for i, item := range raw {
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

func (cli *Cli) GetDirectoryUnit(name string, shared bool) (map[string]any, error) {
	if name == "" {
		return nil, fmt.Errorf("unit name is required")
	}

	const cmd = "GETDIRECTORYUNIT"
	var res any
	var err error

	if shared {
		res, err = cli.Query(cmd, name, "SHARED")
	} else {
		res, err = cli.Query(cmd, name)
	}
	if err != nil {
		return nil, err
	}

	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

func (cli *Cli) ListDirectoryUnits(shared bool) (map[string]string, error) {
	const cmd = "LISTDIRECTORYUNITS"
	var res any
	var err error

	if shared {
		res, err = cli.Query(cmd, "SHARED")
	} else {
		res, err = cli.Query(cmd)
	}
	if err != nil {
		return nil, err
	}

	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}

	// Инвертируем мапу: CGP отдает {Mount: Name}, а нам нужно {Name: Mount}
	units := make(map[string]string, len(m))
	for mnt, nameAny := range m {
		name, _ := toString(cmd, nameAny)
		units[name] = mnt
	}
	return units, nil
}

func (cli *Cli) RelocateDirectoryUnit(name, newMountPoint string, shared bool) (err error) {
	if name == "" || newMountPoint == "" {
		return fmt.Errorf("unit name and new mount point are required")
	}
	const cmd = "RELOCATEDIRECTORYUNIT"
	if shared {
		err = cli.QueryNV(cmd, name, "SHARED", newMountPoint)
	} else {
		err = cli.QueryNV(cmd, name, newMountPoint)
	}
	return
}

func (cli *Cli) SetDirectoryAccessRights(shared bool, rights []DirAccessRight) (err error) {
	const cmd = "SETDIRECTORYACCESSRIGHTS"
	raw := make([]any, len(rights))
	for i, r := range rights {
		raw[i] = []any{
			r.Name, r.Target, r.BindDN, r.Type,
			r.RecordRights, r.SearchAttrs, r.ReadAttrs, r.ModifyAttrs,
		}
	}

	if shared {
		err = cli.QueryNV(cmd, "SHARED", raw)
	} else {
		err = cli.QueryNV(cmd, raw)
	}
	return
}

func (cli *Cli) SetDirectoryUnit(name string, shared bool, settings map[string]any) (err error) {
	if name == "" {
		return fmt.Errorf("unit name is required")
	}
	const cmd = "SETDIRECTORYUNIT"
	if shared {
		err = cli.QueryNV(cmd, name, "SHARED", settings)
	} else {
		err = cli.QueryNV(cmd, name, settings)
	}
	return
}
