// # Section: Web Skins Administration
package cgpcli

import (
	"fmt"
	"time"
)

// CreateClusterSkin создает кастомный кластерный скин.
func (cli *Cli) CreateClusterSkin(skin string) (err error) {
	_, err = cli.Query("CREATECLUSTERSKIN", skin)
	return
}

// CreateDomainSkin создает кастомный скин домена.
func (cli *Cli) CreateDomainSkin(domain, skin string) (err error) {
	const cmd = "CREATEDOMAINSKIN"
	if domain != "" {
		_, err = cli.Query(cmd, domain, "SKIN", skin)
	} else {
		_, err = cli.Query(cmd, skin)
	}
	return
}

// CreateServerSkin создает кастомный серверный скин.
func (cli *Cli) CreateServerSkin(skin string) (err error) {
	_, err = cli.Query("CREATESERVERSKIN", skin)
	return
}

// DeleteClusterSkin удаляет кастомный кластерный скин.
func (cli *Cli) DeleteClusterSkin(skin string) (err error) {
	_, err = cli.Query("DELETECLUSTERSKIN", skin)
	return
}

// DeleteDomainSkin удаляет кастомный скин домена.
func (cli *Cli) DeleteDomainSkin(domain, skin string) (err error) {
	const cmd = "DELETEDOMAINSKIN"
	if domain != "" {
		_, err = cli.Query(cmd, domain, "SKIN", skin)
	} else {
		_, err = cli.Query(cmd, skin)
	}
	return
}

// DeleteServerSkin удаляет кастомный серверный скин.
func (cli *Cli) DeleteServerSkin(skin string) (err error) {
	_, err = cli.Query("DELETESERVERSKIN", skin)
	return
}

// ListClusterSkinFiles возвращает словарь с атрибутами файлов в кластерном скине.
func (cli *Cli) ListClusterSkinFiles(skin string) (map[string]any, error) {
	const cmd = "LISTCLUSTERSKINFILES"
	res, err := cli.Query(cmd, skin)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// ListClusterSkins возвращает список имен кастомных кластерных скинов.
func (cli *Cli) ListClusterSkins() ([]string, error) {
	const cmd = "LISTCLUSTERSKINS"
	res, err := cli.Query(cmd)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ListDomainSkinFiles возвращает словарь с атрибутами файлов в скине домена.
func (cli *Cli) ListDomainSkinFiles(domain, skin string) (map[string]any, error) {
	const cmd = "LISTDOMAINSKINFILES"
	var res any
	var err error

	if domain != "" {
		res, err = cli.Query(cmd, domain, "SKIN", skin)
	} else {
		res, err = cli.Query(cmd, skin)
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

// ListDomainSkins возвращает список имен кастомных скинов домена.
func (cli *Cli) ListDomainSkins(domain string) ([]string, error) {
	const cmd = "LISTDOMAINSKINS"
	var res any
	var err error

	if domain != "" {
		res, err = cli.Query(cmd, domain)
	} else {
		res, err = cli.Query(cmd)
	}

	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ListServerSkinFiles возвращает словарь с атрибутами файлов в серверном скине.
func (cli *Cli) ListServerSkinFiles(skin string) (map[string]any, error) {
	const cmd = "LISTSERVERSKINFILES"
	res, err := cli.Query(cmd, skin)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// ListServerSkins возвращает список имен кастомных серверных скинов.
func (cli *Cli) ListServerSkins() ([]string, error) {
	const cmd = "LISTSERVERSKINS"
	res, err := cli.Query(cmd)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ListStockSkinFiles возвращает список файлов встроенного (stock) скина.
func (cli *Cli) ListStockSkinFiles(skin string) (map[string]any, error) {
	const cmd = "LISTSTOCKSKINFILES"
	res, err := cli.Query(cmd, skin)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

// ReadClusterSkinFile читает файл из кастомного кластерного скина.
func (cli *Cli) ReadClusterSkinFile(skin, file string) ([]byte, time.Time, error) {
	const cmd = "READCLUSTERSKINFILE"
	res, err := cli.Query(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadDomainSkinFile читает файл из кастомного скина домена.
func (cli *Cli) ReadDomainSkinFile(domain, skin, file string) ([]byte, time.Time, error) {
	const cmd = "READDOMAINSKINFILE"
	var res any
	var err error

	if domain != "" {
		res, err = cli.Query(cmd, domain, "SKIN", skin, "FILE", file)
	} else {
		res, err = cli.Query(cmd, skin, "FILE", file)
	}

	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array [data, time], got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadServerSkinFile читает файл из кастомного серверного скина.
func (cli *Cli) ReadServerSkinFile(skin, file string) ([]byte, time.Time, error) {
	const cmd = "READSERVERSKINFILE"
	res, err := cli.Query(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// ReadStockSkinFile читает файл из встроенного (stock) скина.
func (cli *Cli) ReadStockSkinFile(skin, file string) ([]byte, time.Time, error) {
	const cmd = "READSTOCKSKINFILE"
	res, err := cli.Query(cmd, skin, "FILE", file)
	if err != nil {
		return nil, time.Time{}, err
	}

	arr, ok := res.([]any)
	if !ok || len(arr) < 2 {
		return nil, time.Time{}, fmt.Errorf("%s: expected array, got %T", cmd, res)
	}

	data, ok1 := arr[0].([]byte)
	ts, ok2 := arr[1].(time.Time)
	if !ok1 || !ok2 {
		return nil, time.Time{}, fmt.Errorf("%s: invalid data types in response", cmd)
	}

	return data, ts, nil
}

// RenameClusterSkin переименовывает кастомный кластерный скин.
func (cli *Cli) RenameClusterSkin(skin, newSkinID string) (err error) {
	_, err = cli.Query("RENAMECLUSTERSKIN", skin, "INTO", newSkinID)
	return
}

// RenameDomainSkin переименовывает кастомный скин домена.
func (cli *Cli) RenameDomainSkin(domain, skin, newSkinID string) (err error) {
	const cmd = "RENAMEDOMAINSKIN"
	if domain != "" {
		_, err = cli.Query(cmd, domain, "SKIN", skin, "INTO", newSkinID)
	} else {
		_, err = cli.Query(cmd, skin, "INTO", newSkinID)
	}
	return
}

// RenameServerSkin переименовывает кастомный серверный скин.
func (cli *Cli) RenameServerSkin(skin, newSkinID string) (err error) {
	_, err = cli.Query("RENAMESERVERSKIN", skin, "INTO", newSkinID)
	return
}

// StoreClusterSkinFile сохраняет файл в кастомный кластерный скин или удаляет его.
func (cli *Cli) StoreClusterSkinFile(skin, file string, data []byte) (err error) {
	const cmd = "STORECLUSTERSKINFILE"
	if data != nil {
		_, err = cli.Query(cmd, skin, "FILE", file, "DATA", data)
	} else {
		_, err = cli.Query(cmd, skin, "FILE", file, "DELETE")
	}
	return
}

// StoreDomainSkinFile сохраняет файл в кастомный скин домена или удаляет его.
func (cli *Cli) StoreDomainSkinFile(domain, skin, file string, data []byte) (err error) {
	args := make([]any, 0, 8)
	args = append(args, "STOREDOMAINSKINFILE")
	if domain != "" {
		args = append(args, domain, "SKIN")
	}
	args = append(args, skin, "FILE", file)

	if data != nil {
		args = append(args, "DATA", data)
	} else {
		args = append(args, "DELETE")
	}

	_, err = cli.Query(args...)
	return
}

// StoreServerSkinFile сохраняет файл в кастомный серверный скин или удаляет его.
func (cli *Cli) StoreServerSkinFile(skin, file string, data []byte) (err error) {
	const cmd = "STORESERVERSKINFILE"
	if data != nil {
		_, err = cli.Query(cmd, skin, "FILE", file, "DATA", data)
	} else {
		_, err = cli.Query(cmd, skin, "FILE", file, "DELETE")
	}
	return
}
