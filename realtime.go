// # Section: Real-Time Application Administration
package cgpcli

import (
	"fmt"
)

func (cli *Cli) CreateClusterPBX(language string) (err error) {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	_, err = cli.Query("CREATECLUSTERPBX", language)
	return
}

// CreateDomainPBX создает среду RTA для домена или её национальное подмножество.
func (cli *Cli) CreateDomainPBX(domain, language string) (err error) {
	if domain == "" {
		return fmt.Errorf("domain is required")
	}
	const cmd = "CREATEDOMAINPBX"
	if language != "" {
		_, err = cli.Query(cmd, domain, "FILE", language)
	} else {
		_, err = cli.Query(cmd, domain)
	}
	return
}

func (cli *Cli) CreateServerPBX(language string) (err error) {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	_, err = cli.Query("CREATESERVERPBX", language)
	return
}

func (cli *Cli) DeleteClusterPBX(language string) (err error) {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	_, err = cli.Query("DELETECLUSTERPBX", language)
	return
}

// DeleteDomainPBX удаляет национальное подмножество среды RTA домена.
func (cli *Cli) DeleteDomainPBX(domain, language string) (err error) {
	if domain == "" || language == "" {
		return fmt.Errorf("domain and language are required")
	}
	_, err = cli.Query("DELETEDOMAINPBX", domain, "FILE", language)
	return
}

func (cli *Cli) DeleteServerPBX(language string) (err error) {
	if language == "" {
		return fmt.Errorf("language is required")
	}
	_, err = cli.Query("DELETESERVERPBX", language)
	return
}

func (cli *Cli) ListClusterPBXFiles(language string) (map[string]any, error) {
	const cmd = "LISTCLUSTERPBXFILES"
	var res any
	var err error
	if language != "" {
		res, err = cli.Query(cmd, language)
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
	return m, nil
}

// ListDomainPBXFiles возвращает список файлов в среде RTA домена.
func (cli *Cli) ListDomainPBXFiles(domain, language string) (map[string]any, error) {
	const cmd = "LISTDOMAINPBXFILES"
	args := make([]any, 0, 3)
	if domain != "" {
		args = append(args, domain)
	}
	if language != "" {
		args = append(args, "FILE", language)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	m, ok := res.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected map, got %T", cmd, res)
	}
	return m, nil
}

func (cli *Cli) ListServerPBXFiles(language string) (map[string]any, error) {
	const cmd = "LISTSERVERPBXFILES"
	var res any
	var err error
	if language != "" {
		res, err = cli.Query(cmd, language)
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
	return m, nil
}

func (cli *Cli) ListStockPBXFiles(language string) (map[string]any, error) {
	const cmd = "LISTSTOCKPBXFILES"
	var res any
	var err error
	if language != "" {
		res, err = cli.Query(cmd, language)
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
	return m, nil
}

func (cli *Cli) ReadClusterPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	const cmd = "READCLUSTERPBXFILE"
	res, err := cli.Query(cmd, file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
}

// ReadDomainPBXFile читает файл из среды RTA домена.
func (cli *Cli) ReadDomainPBXFile(domain, file string) ([]byte, error) {
	if domain == "" || file == "" {
		return nil, fmt.Errorf("domain and file are required")
	}
	const cmd = "READDOMAINPBXFILE"
	res, err := cli.Query(cmd, domain, "FILE", file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
}

func (cli *Cli) ReadServerPBXFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("file is required")
	}
	const cmd = "READSERVERPBXFILE"
	res, err := cli.Query(cmd, file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
}

func (cli *Cli) ReadStockPBXFile(file string) ([]byte, error) {
	const cmd = "READSTOCKPBXFILE"
	res, err := cli.Query(cmd, file)
	if err != nil {
		return nil, err
	}
	data, ok := res.([]byte)
	if !ok {
		return nil, fmt.Errorf("%s: expected []byte, got %T", cmd, res)
	}
	return data, nil
}

func (cli *Cli) StoreClusterPBXFile(file string, data []byte) (err error) {
	if file == "" {
		return fmt.Errorf("file is required")
	}
	const cmd = "STORECLUSTERPBXFILE"
	if data == nil {
		_, err = cli.Query(cmd, file, "DELETE")
	} else {
		_, err = cli.Query(cmd, file, "DATA", data)
	}
	return
}

// StoreDomainPBXFile сохраняет или удаляет файл в среде RTA домена.
func (cli *Cli) StoreDomainPBXFile(domain, file string, data []byte) (err error) {
	if domain == "" || file == "" {
		return fmt.Errorf("domain and file are required")
	}
	const cmd = "STOREDOMAINPBXFILE"
	if data == nil {
		_, err = cli.Query(cmd, domain, "FILE", file, "DELETE")
	} else {
		_, err = cli.Query(cmd, domain, "FILE", file, "DATA", data)
	}
	return
}

func (cli *Cli) StoreServerPBXFile(file string, data []byte) (err error) {
	if file == "" {
		return fmt.Errorf("file is required")
	}
	const cmd = "STORESERVERPBXFILE"
	if data == nil {
		_, err = cli.Query(cmd, file, "DELETE")
	} else {
		_, err = cli.Query(cmd, file, "DATA", data)
	}
	return
}
