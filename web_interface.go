// # Web Interface Integration
//
// The Web Interface section provides methods for managing user sessions across
// various CommuniGate Pro protocols, including WebUser, XIMSS, and Lite sessions.
// These methods are essential for building custom web portals, implementing
// Single Sign-On (SSO), and managing active connections.
//
// Key capabilities include:
//   - Session Creation: programmatically initiating WebUser, XIMSS, or Lite
//     sessions for accounts via [Cli.CreateWebUserSession] and related methods.
//   - Multi-Factor Authentication: supporting the 2FA flow by finalizing
//     authentication steps with [Cli.BlessSession].
//   - Session Discovery & Inspection: locating active sessions by IP address,
//     protocol, or client type using [Cli.FindAccountSession] and [Cli.GetSession].
//   - Resource Management: handling session-bound file uploads and storage
//     via [Cli.StoreSessionFile], supporting various offsets for partial data.
//   - Lifecycle Control: terminating specific sessions with [Cli.KillSession]
//     or updating session-specific metadata with [Cli.UpdateSession].
package cgpcli

// SessionOffset specifies the file position or write mode for session file operations.
//
// Values:
//   - (positive number): the File Storage file is rewritten/extended starting
//     from the specified file position. The file should already exist.
//   - OffsetBeg: the file is rewritten from the beginning, but its old data
//     beyond the end of the new data is not removed.
//   - OffsetEnd: the file data is appended to the end of the file. If the
//     file does not exist, it is created.
//   - OffsetNew: a new file is created and file data is stored in it. The
//     file must not exist.
//   - (zero/absent): the existing file is removed first, and a new file is created.
type SessionOffset int

const (
	OffsetBeg SessionOffset = -1
	OffsetEnd SessionOffset = -2
	OffsetNew SessionOffset = -3
)

// BlessSession completes the second stage of a Two-factor authentication
// process for the given session.
//
// Parameters:
//   - sessionID: the Session ID.
//   - secret: an optional one-time secret used with Two-factor
//     Authentication.
//   - account: an optional Account name the session belongs to.
//
// This method executes the BLESSSESSION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) BlessSession(sessionID string, secret, account string) error {
	if sessionID == "" {
		return ErrSessionIDRequired
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

// CreateLiteSession creates a LITE session.
//
// Parameters:
//   - addr: the IP address and port of the client browser.
//   - origAddr: an optional original IP address of the client
//     browser, if the client connects via a proxy.
//
// This method executes the CREATELITESESSION CLI command.
//
// Returns:
//   - string: the LITE Session ID.
//   - error: an error if the command fails.
func (cli *Cli) CreateLiteSession(addr, origAddr string) (string, error) {
	if addr == "" {
		return "", ErrAddrRequired
	}

	const cmd = "CREATELITESESSION"
	if origAddr != "" {
		return cli.getString(cmd, "ADDRESS", addr, "FOR", origAddr)
	} else {
		return cli.getString(cmd, "ADDRESS", addr)
	}
}

// CreateWebUserSession creates a WebUser session for the specified Account.
//
// Parameters:
//   - account: the Account name.
//   - addr: the IP address and port of the client browser.
//   - origAddr: an optional original IP address of the client
//     browser, if the client connects via a proxy.
//   - skin: an optional Skin name to use for the newly created session.
//
// This method executes the CREATEWEBUSERSESSION CLI command.
//
// Returns:
//   - string: the WebUser Session ID.
//   - error: an error if the command fails.
func (cli *Cli) CreateWebUserSession(account, addr, origAddr, skin string) (string, error) {
	if account == "" {
		return "", ErrAccountRequired
	}
	if addr == "" {
		return "", ErrAddrRequired
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

// CreateXimssSession creates a XIMSS session for the specified Account.
//
// Parameters:
//   - account: the Account name.
//   - addr: the IP address and port of the client browser.
//   - origAddr: an optional original IP address of the client
//     browser, if the client connects via a proxy.
//
// This method executes the CREATEXIMSSSESSION CLI command.
//
// Returns:
//   - string: the XIMSS Session ID.
//   - error: an error if the command fails.
func (cli *Cli) CreateXimssSession(account, addr, origAddr string) (string, error) {
	if account == "" {
		return "", ErrAccountRequired
	}
	if addr == "" {
		return "", ErrAddrRequired
	}

	const cmd = "CREATEXIMSSSESSION"
	if origAddr != "" {
		return cli.getString(cmd, account, "ADDRESS", addr, "FOR", origAddr)
	} else {
		return cli.getString(cmd, account, "ADDRESS", addr)
	}
}

// FindAccountSession finds an existing session for the specified Account.
//
// Parameters:
//   - account: the Account name.
//   - addr: an optional IP address of the client browser.
//   - origAddr: an optional IP address of the client browser,
//     if this browser is located behind an HTTP proxy.
//   - protocol: an optional Session protocol (WebUser, XIMSS, XMPP, etc.).
//   - transport: an optional Session transport (HTTP, XIMSS, XMPP, etc.).
//   - client: an optional Session client.
//
// This method executes the ADDRESS CLI command.
//
// Returns:
//   - string: the Session ID.
//   - error: an error if the command fails.
func (cli *Cli) FindAccountSession(account, addr, origAddr, protocol, transport, client string) (string, error) {
	if account == "" {
		return "", ErrAccountRequired
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

// GetAccountWebAuthSession retrieves the Web Auth session data for the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - auth: the Web Auth Session ID.
//
// This data contains the session dataset for Web Auth authentication.
//
// This method executes the GETACCOUNTWEBAUTHSESSION CLI command.
//
// Returns:
//   - map[string]any: a map containing the session dataset.
//   - error: an error if the command fails.
func (cli *Cli) GetAccountWebAuthSession(account, auth string) (map[string]any, error) {
	if account == "" {
		return nil, ErrAccountRequired
	}
	if auth == "" {
		return nil, ErrAuthRequired
	}
	return cli.getMapAny("GETACCOUNTWEBAUTHSESSION", account, "ID", auth)
}

// GetSession retrieves Session data.
//
// Parameters:
//   - sessionID: the Session ID.
//   - domain: an optional Domain name the session Account belongs to.
//
// This method executes the GETSESSION CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the session dataset.
//   - error: an error if the command fails.
func (cli *Cli) GetSession(sessionID, domain string) (map[string]any, error) {
	if sessionID == "" {
		return nil, ErrSessionIDRequired
	}

	const cmd = "GETSESSION"
	if domain != "" {
		return cli.getMapAny(cmd, sessionID, "DOMAIN", domain)
	} else {
		return cli.getMapAny(cmd, sessionID)
	}
}

// KillAccountWebAuthSession terminate a Web Auth session for the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - auth: the Web Auth Session ID.
//
// Note: This command also terminates related XIMSS sessions.
//
// This method executes the KILLACCOUNTWEBAUTHSESSION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) KillAccountWebAuthSession(account, auth string) (int64, error) {
	if account == "" {
		return 0, ErrAccountRequired
	}
	if auth == "" {
		return 0, ErrAuthRequired
	}
	return cli.getInt("KILLACCOUNTWEBAUTHSESSION", account, "ID", auth)
}

// KillSession terminates a Session.
//
// Parameters:
//   - sessionID: the Session ID.
//   - domain: an optional Domain name the session Account belongs to.
//
// This method executes the KILLSESSION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) KillSession(sessionID, domain string) error {
	if sessionID == "" {
		return ErrSessionIDRequired
	}

	if domain != "" {
		return cli.QueryNV("KILLSESSION", sessionID, "DOMAIN", domain)
	} else {
		return cli.QueryNV("KILLSESSION", sessionID)
	}
}

// ListAccountSessions retrieves all existing sessions for the specified Account.
//
// Parameters:
//   - account: the Account name.
//   - addr: an optional IP address of the client browser.
//   - origAddr: an optional IP address of the client browser,
//     if this browser is located behind an HTTP proxy.
//   - protocol: an optional Session protocol (WebUser, XIMSS, XMPP, etc.).
//   - transport: an optional Session transport (HTTP, XIMSS, XMPP, etc.).
//   - client: an optional Session client.
//
// This method executes the ADDRESS CLI command.
//
// Returns:
//   - []string: an array of Session IDs.
//   - error: an error if the command fails.
func (cli *Cli) ListAccountSessions(account, addr, origAddr, protocol, transport, client string) ([]string, error) {
	if account == "" {
		return nil, ErrAccountRequired
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

// ListAccountWebAuthSessions retrieve all existing Web Auth sessions for the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - activeonly: only WebAuth session IDs that have active sessions are returned.
//
// This method executes the LISTACCOUNTWEBAUTHSESSIONS CLI command.
//
// Returns:
//   - []string: an array of strings, where each string is the WebAuth Session ID.
//   - error: an error if the command fails.
func (cli *Cli) ListAccountWebAuthSessions(account string, activeonly bool) ([]string, error) {
	if account == "" {
		return nil, ErrAccountRequired
	}
	if activeonly {
		return cli.getSliceString("LISTACCOUNTWEBAUTHSESSIONS", account, "ACTIVEONLY")
	} else {
		return cli.getSliceString("LISTACCOUNTWEBAUTHSESSIONS", account)
	}
}

// ListLiteSessions retrieves all LITE sessions.
//
// Parameters:
//   - addr: an optional IP address of the client browser.
//   - origAddr: an optional IP address of the client browser,
//     if this browser is located behind an HTTP proxy.
//
// This method executes the ADDRESS CLI command.
//
// Returns:
//   - []string: an array of Session IDs.
//   - error: an error if the command fails.
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

// StoreSessionFile stores an uploaded file from the session "uploaded file set"
// as a File Storage file.
//
// Parameters:
//   - sessionID: the Session ID.
//   - domain: an optional Domain name the session Account belongs to.
//   - file: the file name.
//   - uploadID: identifies a file in the "uploaded file set".
//   - offset: an optional file position. Can be a positive number,
//     or constants OffsetBeg, OffsetEnd, OffsetNew.
//
// This method executes the DOMAIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) StoreSessionFile(sessionID, domain, file, uploadID string, offset SessionOffset) error {
	if sessionID == "" {
		return ErrSessionIDRequired
	}
	if file == "" {
		return ErrFileRequired
	}
	if uploadID == "" {
		return ErrUploadIDRequired
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

// UpdateSession modifies custom parameters of a Session.
//
// Parameters:
//   - sessionID: the Session ID.
//   - domain: an optional Domain name the session Account belongs to.
//   - param: a dictionary of Session attribute names and their new values.
//     Set a value to nil to remove the attribute.
//
// This method executes the UPDATESESSION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateSession(sessionID, domain string, param map[string]any) error {
	if domain != "" {
		return cli.QueryNV("UPDATESESSION", sessionID, "DOMAIN", domain, param)
	} else {
		return cli.QueryNV("UPDATESESSION", sessionID, param)
	}
}
