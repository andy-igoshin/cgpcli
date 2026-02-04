// # Section: Web Interface Integration
package cgpcli

import (
	"fmt"
)

// BlessSession завершает второй этап двухфакторной аутентификации (2FA).
func (cli *Cli) BlessSession(sessionID string, secret, account string) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID is required")
	}

	args := make([]any, 0, 5)
	args = append(args, sessionID)
	if secret != "" {
		args = append(args, "PASSWORD", secret)
	}
	if account != "" {
		args = append(args, "AUTH", account)
	}
	return cli.QueryNV("BLESSSESSION", args...)
}

// CreateLiteSession создает новую LITE-сессию.
func (cli *Cli) CreateLiteSession(addr, origAddr string) (s string, err error) {
	if addr == "" {
		return "", fmt.Errorf("addr is required")
	}

	const cmd = "CREATELITESESSION"
	if origAddr != "" {
		s, err = cli.getString(cmd, "ADDRESS", addr, "FOR", origAddr)
	} else {
		s, err = cli.getString(cmd, "ADDRESS", addr)
	}
	return
}

// CreateWebUserSession создает сессию WebUser для указанного аккаунта.
func (cli *Cli) CreateWebUserSession(account, addr, origAddr, skin string) (string, error) {
	if account == "" || addr == "" {
		return "", fmt.Errorf("account name and addr are required")
	}

	const cmd = "CREATEWEBUSERSESSION"
	args := make([]any, 0, 7)
	args = append(args, account, "ADDRESS", addr)
	if origAddr != "" {
		args = append(args, "FOR", origAddr)
	}
	if skin != "" {
		args = append(args, "SKIN", skin)
	}
	return cli.getString(cmd, args...)
}

// CreateXimssSession создает сессию XIMSS для указанного аккаунта.
func (cli *Cli) CreateXimssSession(account, addr, origAddr string) (s string, err error) {
	if account == "" || addr == "" {
		return "", fmt.Errorf("account name and addr are required")
	}

	const cmd = "CREATEXIMSSSESSION"
	if origAddr != "" {
		s, err = cli.getString(cmd, account, "ADDRESS", addr, "FOR", origAddr)
	} else {
		s, err = cli.getString(cmd, account, "ADDRESS", addr)
	}
	return
}

// FindAccountSession ищет существующую сессию аккаунта по заданным критериям.
func (cli *Cli) FindAccountSession(account, addr, origAddr, protocol, transport, client string) (string, error) {
	if account == "" {
		return "", fmt.Errorf("account name is required")
	}

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
	return cli.getString("FINDACCOUNTSESSION", args...)
}

// GetSession возвращает данные сессии в виде словаря.
func (cli *Cli) GetSession(sessionID, domain string) (res map[string]any, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("sessionID is required")
	}

	const cmd = "GETSESSION"
	if domain != "" {
		res, err = cli.getMapAny(cmd, sessionID, "DOMAIN", domain)
	} else {
		res, err = cli.getMapAny(cmd, sessionID)
	}
	return
}

// KillSession принудительно завершает сессию.
func (cli *Cli) KillSession(sessionID, domain string) (err error) {
	if sessionID == "" {
		return fmt.Errorf("sessionID is required")
	}

	if domain != "" {
		err = cli.QueryNV("KILLSESSION", sessionID, "DOMAIN", domain)
	} else {
		err = cli.QueryNV("KILLSESSION", sessionID)
	}
	return
}

// ListAccountSessions возвращает список ID всех существующих сессий для аккаунта.
func (cli *Cli) ListAccountSessions(account, addr, origAddr, protocol, transport, client string) ([]string, error) {
	if account == "" {
		return nil, fmt.Errorf("account name is required")
	}

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
	return cli.getSliceString("LISTACCOUNTSESSIONS", args...)
}

// ListLiteSessions возвращает список всех LITE-сессий.
func (cli *Cli) ListLiteSessions(addr, origAddr string) ([]string, error) {
	args := make([]any, 0, 4)
	if addr != "" {
		args = append(args, "ADDRESS", addr)
		if origAddr != "" {
			args = append(args, "FOR", origAddr)
		}
	}
	return cli.getSliceString("LISTLITESESSIONS", args...)
}

// StoreSessionFile сохраняет файл из сессионного набора в File Storage.
// offset: может быть позицией (>=0) или константами OffsetBeg, OffsetEnd, OffsetNew.
func (cli *Cli) StoreSessionFile(sessionID, domain, file, uploadID string, offset SessionOffset) error {
	if sessionID == "" || file == "" || uploadID == "" {
		return fmt.Errorf("sessionID, file name and uploadID are required")
	}

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
	return cli.QueryNV("STORESESSIONFILE", args...)
}

// UpdateSession изменяет кастомные параметры сессии.
func (cli *Cli) UpdateSession(sessionID, domain string, data map[string]any) (err error) {
	if domain != "" {
		err = cli.QueryNV("UPDATESESSION", sessionID, "DOMAIN", domain, data)
	} else {
		err = cli.QueryNV("UPDATESESSION", sessionID, data)
	}
	return
}
