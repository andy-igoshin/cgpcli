// # Section: Web Interface Integration
package cgpcli

import (
	"fmt"
)

// BlessSession завершает второй этап двухфакторной аутентификации (2FA).
func (cli *Cli) BlessSession(sessionID string, secret, account string) (err error) {
	args := make([]any, 0, 5)
	args = append(args, sessionID)
	if secret != "" {
		args = append(args, "PASSWORD", secret)
	}
	if account != "" {
		args = append(args, "AUTH", account)
	}
	_, err = cli.Query("BLESSSESSION", args...)
	return
}

// CreateLiteSession создает новую LITE-сессию.
func (cli *Cli) CreateLiteSession(addr, origAddr string) (string, error) {
	const cmd = "CREATELITESESSION"
	var res any
	var err error

	if origAddr != "" {
		res, err = cli.Query(cmd, "ADDRESS", addr, "FOR", origAddr)
	} else {
		res, err = cli.Query(cmd, "ADDRESS", addr)
	}

	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
	}
	return s, nil
}

// CreateWebUserSession создает сессию WebUser для указанного аккаунта.
func (cli *Cli) CreateWebUserSession(account, addr, origAddr, skin string) (string, error) {
	const cmd = "CREATEWEBUSERSESSION"
	args := make([]any, 0, 7)
	args = append(args, account, "ADDRESS", addr)
	if origAddr != "" {
		args = append(args, "FOR", origAddr)
	}
	if skin != "" {
		args = append(args, "SKIN", skin)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
	}
	return s, nil
}

// CreateXimssSession создает сессию XIMSS для указанного аккаунта.
func (cli *Cli) CreateXimssSession(account, addr, origAddr string) (string, error) {
	const cmd = "CREATEXIMSSSESSION"
	var res any
	var err error

	if origAddr != "" {
		res, err = cli.Query(cmd, account, "ADDRESS", addr, "FOR", origAddr)
	} else {
		res, err = cli.Query(cmd, account, "ADDRESS", addr)
	}

	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
	}
	return s, nil
}

// FindAccountSession ищет существующую сессию аккаунта по заданным критериям.
func (cli *Cli) FindAccountSession(account, addr, origAddr, protocol, transport, client string) (string, error) {
	const cmd = "FINDACCOUNTSESSION"
	args := make([]any, 0, 11)
	args = append(args, account)
	if addr != "" {
		args = append(args, "ADDRESS", addr)
		if origAddr != "" {
			args = append(args, "FOR", origAddr)
		}
	}
	if protocol != "" {
		args = append(args, "PROTOCOL", protocol)
	}
	if transport != "" {
		args = append(args, "TRANSPORT", transport)
	}
	if client != "" {
		args = append(args, "CLIENT", client)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return "", err
	}
	s, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("%s: expected string, got %T", cmd, res)
	}
	return s, nil
}

// GetSession возвращает данные сессии в виде словаря.
func (cli *Cli) GetSession(sessionID, domain string) (map[string]any, error) {
	const cmd = "GETSESSION"
	var res any
	var err error

	if domain == "" {
		res, err = cli.Query(cmd, sessionID)
	} else {
		res, err = cli.Query(cmd, sessionID, "DOMAIN", domain)
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

// KillSession принудительно завершает сессию.
func (cli *Cli) KillSession(sessionID, domain string) (err error) {
	if domain == "" {
		_, err = cli.Query("KILLSESSION", sessionID)
	} else {
		_, err = cli.Query("KILLSESSION", sessionID, "DOMAIN", domain)
	}
	return
}

// ListAccountSessions возвращает список ID всех существующих сессий для аккаунта.
func (cli *Cli) ListAccountSessions(account, addr, origAddr, protocol, transport, client string) ([]string, error) {
	const cmd = "LISTACCOUNTSESSIONS"
	args := make([]any, 0, 11)
	args = append(args, account)
	if addr != "" {
		args = append(args, "ADDRESS", addr)
		if origAddr != "" {
			args = append(args, "FOR", origAddr)
		}
	}
	if protocol != "" {
		args = append(args, "PROTOCOL", protocol)
	}
	if transport != "" {
		args = append(args, "TRANSPORT", transport)
	}
	if client != "" {
		args = append(args, "CLIENT", client)
	}

	res, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, res)
}

// ListLiteSessions возвращает список всех LITE-сессий.
func (cli *Cli) ListLiteSessions(addr, origAddr string) ([]string, error) {
	const cmd = "LISTLITESESSIONS"
	args := make([]any, 0, 4)

	if addr != "" {
		args = append(args, "ADDRESS", addr)
		if origAddr != "" {
			args = append(args, "FOR", origAddr)
		}
	}

	raw, err := cli.Query(cmd, args...)
	if err != nil {
		return nil, err
	}
	return toStringSlice(cmd, raw)
}

// StoreSessionFile сохраняет файл из сессионного набора в File Storage.
// offset: может быть позицией (>=0) или константами OffsetBeg, OffsetEnd, OffsetNew.
func (cli *Cli) StoreSessionFile(sessionID, domain, file, uploadID string, offset SessionOffset) (err error) {
	args := make([]any, 0, 11)
	args = append(args, sessionID)
	if domain != "" {
		args = append(args, "DOMAIN", domain)
	}
	args = append(args, "FILE", file)

	if offset != 0 {
		args = append(args, "OFFSET")
		switch offset {
		case OffsetBeg:
			args = append(args, "BEG")
		case OffsetEnd:
			args = append(args, "END")
		case OffsetNew:
			args = append(args, "NEW")
		default:
			args = append(args, int(offset))
		}
	}

	args = append(args, "UPLOADID", uploadID)
	_, err = cli.Query("STORESESSIONFILE", args...)
	return
}

// UpdateSession изменяет кастомные параметры сессии.
func (cli *Cli) UpdateSession(sessionID, domain string, data map[string]any) (err error) {
	if domain == "" {
		_, err = cli.Query("UPDATESESSION", sessionID, data)
	} else {
		_, err = cli.Query("UPDATESESSION", sessionID, "DOMAIN", domain, data)
	}
	return
}
