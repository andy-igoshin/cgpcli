// # Section: Directory Administration
package cgpcli

import (
	"fmt"
)

func (cli *Cli) CreateDirectoryUnit(name, mountPoint string, shared, remote bool) error {
	if name == "" || mountPoint == "" {
		return fmt.Errorf("unit name and mount point are required")
	}

	args := make([]any, 0, 5)
	args = append(args, "CREATEDIRECTORYUNIT", name)
	if shared {
		args = append(args, "SHARED")
	}
	if remote {
		args = append(args, "REMOTE")
	}
	args = append(args, mountPoint)

	_, err := cli.Query(args...)
	return err
}

func (cli *Cli) DeleteDirectoryUnit(name string, shared bool) error {
	if name == "" {
		return fmt.Errorf("unit name is required")
	}
	if shared {
		_, err := cli.Query("DELETEDIRECTORYUNIT", name, "SHARED")
		return err
	}
	_, err := cli.Query("DELETEDIRECTORYUNIT", name)
	return err
}

func (cli *Cli) GetDirectoryAccessRights(shared bool) ([]any, error) {
	var res any
	var err error
	if shared {
		res, err = cli.Query("GETDIRECTORYACCESSRIGHTS", "SHARED")
	} else {
		res, err = cli.Query("GETDIRECTORYACCESSRIGHTS")
	}
	if err != nil {
		return nil, err
	}
	arr, ok := res.([]any)
	if !ok {
		return nil, fmt.Errorf("GETDIRECTORYACCESSRIGHTS: expected []any, got %T", res)
	}
	return arr, nil
}

func (cli *Cli) GetDirectoryUnit(name string, shared bool) (map[string]any, error) {
	if name == "" {
		return nil, fmt.Errorf("unit name is required")
	}
	var res any
	var err error
	if shared {
		res, err = cli.Query("GETDIRECTORYUNIT", name, "SHARED")
	} else {
		res, err = cli.Query("GETDIRECTORYUNIT", name)
	}
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("GETDIRECTORYUNIT: expected map, got %T", res)
	}
	return m, nil
}

func (cli *Cli) ListDirectoryUnits(shared bool) (map[string]string, error) {
	var res any
	var err error
	if shared {
		res, err = cli.Query("LISTDIRECTORYUNITS", "SHARED")
	} else {
		res, err = cli.Query("LISTDIRECTORYUNITS")
	}
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("LISTDIRECTORYUNITS: expected map, got %T", res)
	}
	return cli.toStringMap("LISTDIRECTORYUNITS", m)
}

func (cli *Cli) RelocateDirectoryUnit(name, newMountPoint string, shared bool) error {
	if name == "" || newMountPoint == "" {
		return fmt.Errorf("unit name and new mount point are required")
	}
	if shared {
		_, err := cli.Query("RELOCATEDIRECTORYUNIT", name, "SHARED", newMountPoint)
		return err
	}
	_, err := cli.Query("RELOCATEDIRECTORYUNIT", name, newMountPoint)
	return err
}

func (cli *Cli) SetDirectoryAccessRights(shared bool, rights []any) error {
	if shared {
		_, err := cli.Query("SETDIRECTORYACCESSRIGHTS", "SHARED", rights)
		return err
	}
	_, err := cli.Query("SETDIRECTORYACCESSRIGHTS", rights)
	return err
}

func (cli *Cli) SetDirectoryUnit(name string, shared bool, settings map[string]any) error {
	if name == "" {
		return fmt.Errorf("unit name is required")
	}
	if shared {
		_, err := cli.Query("SETDIRECTORYUNIT", name, "SHARED", settings)
		return err
	}
	_, err := cli.Query("SETDIRECTORYUNIT", name, settings)
	return err
}
