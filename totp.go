// # Two-Factor Authentication (TOTP)
//
// The TOTP section provides dedicated methods for managing Time-Based One-Time Password
// (RFC 6238) secret keys for CommuniGate Pro accounts. It streamlines two-factor
// authentication (2FA) provisioning by handling secure key generation and account state management.
//
// Key capabilities include:
//   - Cryptographically secure 160-bit Base32 secret generation on the stack.
//   - Direct setting and clearing of TOTP secrets per account.
//   - One-step account provisioning for QR-code and authenticator app workflows.
//   - Full compatibility with RFC 6238 and standard authenticator applications.
package cgpcli

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
)

var (
	// Unpadded RFC 4648 Base32 кодировщик для TOTP
	b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)
)

// GenerateTotpSecret creates a new cryptographically secure 160-bit (20-byte)
// random secret encoded in Base32 format without padding (RFC 6238 / RFC 4648).
//
// Returns:
//   - string: the generated Base32 secret key (32 characters).
//   - error: an error if the system's cryptographically secure entropy source
//     (crypto/rand) is unavailable.
func GenerateTotpSecret() (string, error) {
	var buf [20]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("generating TOTP entropy: %w", err)
	}
	return b32NoPadding.EncodeToString(buf[:]), nil
}

// ProvisionTotpSecret generates a new 160-bit Base32 TOTP secret and applies it
// to the specified Account in CommuniGate Pro in a single operation.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// This method executes the SETTOTPSECRET CLI command after generating the key via
// [GenerateTotpSecret]. The returned secret can be directly used to build TOTP URIs
// (e.g., otpauth://) or generate QR codes for authenticator apps.
//
// Returns:
//   - string: the newly generated and provisioned Base32 secret key.
//   - error: an error if entropy generation fails or if the CLI command fails.
func (cli *Cli) ProvisionTotpSecret(account string) (string, error) {
	if account == "" {
		return "", errors.New("account name is required")
	}

	secret, err := GenerateTotpSecret()
	if err != nil {
		return "", err
	}
	if err := cli.SetTotpSecret(account, secret); err != nil {
		return "", err
	}
	return secret, nil
}

// ResetTotpSecret clears the TOTP secret for the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//
// This method executes the RESETTOTPSECRET CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ResetTotpSecret(account string) error {
	if account == "" {
		return errors.New("account name is required")
	}
	return cli.QueryNV("RESETTOTPSECRET", Atom(account))
}

// SetTotpSecret configures or clears the TOTP secret for the specified Account.
//
// Parameters:
//   - account: the Account name.
//     The asterisk (*) symbol can be used to specify the current authenticated Account.
//   - secret: the optional TOTP secret key encoded in Base32 format.
//
// Note: If secret is empty (""), the command omits the SECRET parameter, which
// causes CommuniGate Pro to reset/clear the existing TOTP secret for the Account
// (functionally equivalent to [Cli.ResetTotpSecret]).
//
// This method executes the SETTOTPSECRET CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetTotpSecret(account, secret string) error {
	if account == "" {
		return errors.New("account name is required")
	}

	const cmd = "SETTOTPSECRET"
	if secret != "" {
		return cli.QueryNV(cmd, Atom(account), "SECRET", secret)
	} else {
		return cli.QueryNV(cmd, Atom(account))
	}
}
