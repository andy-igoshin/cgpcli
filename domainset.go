// # Domain Set Administration
//
// The Domain Set section provides functionality for global configuration and
// high-level management of the CommuniGate Pro environment. It handles operations
// that affect the entire server or cluster, rather than specific individual domains.
//
// Key capabilities include:
//   - Domain Lifecycle: creating and deleting domains (including [Directory Domains]),
//     renaming, and managing physical storage mount points.
//   - Cluster-Wide Defaults: managing [Cluster Domain Defaults], [Account Defaults],
//     and [Preferences] that synchronize across all cluster nodes.
//   - Server-Wide Defaults: configuring default settings and preferences used for
//     newly created objects at the single-server level.
//   - Security: managing [Trusted Certificates] (CA lists) for both Cluster and Server levels.
//   - Integration: configuring system-wide [Directory Integration] settings and
//     global telephony number (Telnums).
package cgpcli

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// CreateDirectoryDomain creates a new directory-based Domain.
// This operation is allowed only when the Directory-based Domains are enabled.
//
// Parameters:
//   - domain: the Domain name.
//   - settings: an optional dictionary specifying the Domain settings.
//
// This method executes the CREATEDIRECTORYDOMAIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateDirectoryDomain(domain string, settings map[string]any) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "CREATEDIRECTORYDOMAIN"
	if settings != nil {
		return cli.QueryNV(cmd, domain, settings)
	} else {
		return cli.QueryNV(cmd, domain)
	}
}

// CreateDomain creates a new secondary Domain.
//
// Parameters:
//   - domain: the Domain name.
//   - shared: if true, creates a Cluster-wide Domain (Dynamic Cluster only).
//   - storage: an optional "storage mount point" directory name (without .mnt suffix).
//   - settings: an optional dictionary specifying the initial Domain settings.
//
// This method executes the CREATEDOMAIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateDomain(domain string, settings map[string]any, storage string, shared bool) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	args := make([]any, 0, 5)
	args = append(args, domain)
	if shared {
		args = append(args, "SHARED")
	}
	if storage != "" {
		args = append(args, "PATH", storage)
	}
	if settings != nil {
		args = append(args, settings)
	}
	return cli.QueryNV("CREATEDOMAIN", args...)
}

// CreateDomainStorage creates a "storage mount point" for new Domains.
//
// Parameters:
//   - storage: the "storage mount Point" name.
//   - shared: if true, creates a mount point for Cluster Domains in a Dynamic Cluster.
//
// This method executes the CREATEDOMAINSTORAGE CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) CreateDomainStorage(storage string, shared bool) error {
	if storage == "" {
		return fmt.Errorf("storage path is required")
	}

	const cmd = "CREATEDOMAINSTORAGE"
	if shared {
		return cli.QueryNV(cmd, "SHARED", "PATH", storage)
	} else {
		return cli.QueryNV(cmd, "PATH", storage)
	}
}

// DeleteDomain removes an existing Domain.
//
// Parameters:
//   - domain: the Domain name.
//   - force: if true, removes the Domain even if it contains accounts or other objects.
//
// This method executes the DELETEDOMAIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) DeleteDomain(domain string, force bool) error {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "DELETEDOMAIN"
	if force {
		return cli.QueryNV(cmd, domain, "FORCE")
	} else {
		return cli.QueryNV(cmd, domain)
	}
}

// GetClusterAccountDefaults retrieves the cluster-wide Default Account settings.
//
// This method executes the GETCLUSTERACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the cluster default Account settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterAccountDefaults() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERACCOUNTDEFAULTS")
}

// GetClusterAccountPrefs retrieves the cluster-wide Default Account Preferences.
//
// This method executes the GETCLUSTERACCOUNTPREFS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the cluster default Preferences.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterAccountPrefs() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERACCOUNTPREFS")
}

// GetClusterDirectoryIntegration retrieves the cluster-wide Directory Integration settings.
//
// This method executes the GETCLUSTERDIRECTORYINTEGRATION CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the cluster Directory Integration settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterDirectoryIntegration() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDIRECTORYINTEGRATION")
}

// GetClusterDomainDefaults retrieves the cluster-wide default Domain Settings.
//
// This method executes the GETCLUSTERDOMAINDEFAULTS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the cluster-wide default Domain Settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterDomainDefaults() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDOMAINDEFAULTS")
}

// GetClusterTrustedCerts retrieves the cluster-wide set of Trusted Certificates.
//
// This method executes the GETCLUSTERTRUSTEDCERTS CLI command.
//
// Returns:
//   - []string: a slice of Base64-encoded X.509 certificate data.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterTrustedCerts() ([]string, error) {
	return cli.getTrustedCerts("GETCLUSTERTRUSTEDCERTS")
}

// GetDirectoryIntegration retrieves the server-wide Directory Integration settings.
//
// This method executes the GETDIRECTORYINTEGRATION CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the Directory Integration settings.
//   - error: an error if the command fails.
func (cli *Cli) GetDirectoryIntegration() (map[string]any, error) {
	return cli.getMapAny("GETDIRECTORYINTEGRATION")
}

// GetDomainDefaults retrieves the server-wide default Domain Settings.
//
// This method executes the GETDOMAINDEFAULTS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the default Domain Settings.
//   - error: an error if the command fails.
func (cli *Cli) GetDomainDefaults() (map[string]any, error) {
	return cli.getMapAny("GETDOMAINDEFAULTS")
}

// GetServerAccountDefaults retrieves the server-wide Default Account settings.
//
// This method executes the GETSERVERACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the default Account settings.
//   - error: an error if the command fails.
func (cli *Cli) GetServerAccountDefaults() (map[string]any, error) {
	return cli.getMapAny("GETSERVERACCOUNTDEFAULTS")
}

// GetServerAccountPrefs retrieves the server-wide Default Account Preferences.
//
// This method executes the GETSERVERACCOUNTPREFS CLI command.
//
// Returns:
//   - map[string]any: a dictionary with the default Preferences.
//   - error: an error if the command fails.
func (cli *Cli) GetServerAccountPrefs() (map[string]any, error) {
	return cli.getMapAny("GETSERVERACCOUNTPREFS")
}

// GetServerTrustedCerts retrieves the server-wide set of Trusted Certificates.
//
// This method executes the GETSERVERTRUSTEDCERTS CLI command.
//
// Returns:
//   - []string: an array of Base64-encoded X.509 certificate data.
//   - error: an error if the command fails.
func (cli *Cli) GetServerTrustedCerts() ([]string, error) {
	return cli.getTrustedCerts("GETSERVERTRUSTEDCERTS")
}

// ListClusterTelnums reads Telnum numbers created in shared Cluster Domains.
//
// Parameters:
//   - filter: an optional substring filter for telnum numbers.
//   - limit: maximum number of Telnum numbers to return.
//
// This method executes the LISTCLUSTERTELNUMS CLI command.
//
// Returns:
//   - map[string]string: a dictionary mapping Telnums to Account names.
//   - error: an error if the command fails.
func (cli *Cli) ListClusterTelnums(filter string, limit int) (map[string]string, error) {
	return cli.listTelnums("LISTCLUSTERTELNUMS", filter, limit)
}

// ListDomainStorage retrieves the list of "storage mount points" for Domains.
//
// Parameters:
//   - shared: if true, lists mount points for Cluster Domains.
//
// This method executes the LISTDOMAINSTORAGE CLI command.
//
// Returns:
//   - []string: an array containing the "storage mount points" names.
//   - error: an error if the command fails.
func (cli *Cli) ListDomainStorage(shared bool) ([]string, error) {
	const cmd = "LISTDOMAINSTORAGE"
	if shared {
		return cli.getSliceString(cmd, "SHARED")
	} else {
		return cli.getSliceString(cmd)
	}
}

// ListDomains retrieves the names of all server domains.
//
// This method executes the LISTDOMAINS CLI command.
//
// Returns:
//   - []string: an array with the names of all server domains.
//   - error: an error if the command fails.
func (cli *Cli) ListDomains() ([]string, error) {
	return cli.getSliceString("LISTDOMAINS")
}

// ListServerTelnums reads Telnum numbers created in all non-clustered Domains.
//
// Parameters:
//   - filter: an optional substring filter for telnum numbers.
//   - limit: maximum number of Telnum numbers to return.
//
// This method executes the LISTSERVERTELNUMS CLI command.
//
// Returns:
//   - map[string]string: a dictionary mapping Telnums to Account names.
//   - error: an error if the command fails.
func (cli *Cli) ListServerTelnums(filter string, limit int) (map[string]string, error) {
	return cli.listTelnums("LISTSERVERTELNUMS", filter, limit)
}

// MainDomainName retrieves the name of the Main Domain.
//
// This method executes the MAINDOMAINNAME CLI command.
//
// Returns:
//   - string: the Main Domain name.
//   - error: an error if the command fails.
func (cli *Cli) MainDomainName() (string, error) {
	return cli.getString("MAINDOMAINNAME")
}

// ReloadDirectoryDomains forces the server to scan the Domains Directory subtree.
// This allows finding additional Directory-based Domains created bypassing the CLI.
//
// This method executes the RELOADDIRECTORYDOMAINS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) ReloadDirectoryDomains() error {
	return cli.QueryNV("RELOADDIRECTORYDOMAINS")
}

// RenameDomain changes the name of an existing secondary Domain.
//
// Parameters:
//   - oldName: the name of an existing secondary Domain.
//   - newName: the new name for the Domain.
//   - storage: an optional new "storage mount Point" directory for the Domain data.
//
// This method executes the RENAMEDOMAIN CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RenameDomain(oldName, newName string, storage string) error {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new domain names are required")
	}

	const cmd = "RENAMEDOMAIN"
	if storage != "" {
		return cli.QueryNV(cmd, oldName, "INTO", newName, "PATH", storage)
	} else {
		return cli.QueryNV(cmd, oldName, "INTO", newName)
	}
}

// SetClusterAccountDefaults replaces the cluster-wide Default Account settings.
//
// Parameters:
//   - settings: a dictionary used to replace the Default Account settings.
//
// This method executes the SETCLUSTERACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERACCOUNTDEFAULTS", settings)
}

// SetClusterAccountPrefs replaces the cluster-wide Default Account Preferences.
//
// Parameters:
//   - settings: a dictionary used to replace the Default Account Preferences.
//
// This method executes the SETCLUSTERACCOUNTPREFS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERACCOUNTPREFS", settings)
}

// SetClusterDirectoryIntegration replaces the cluster-wide Directory Integration settings.
//
// Parameters:
//   - settings: a dictionary used to replace the Directory Integration settings.
//
// This method executes the SETCLUSTERDIRECTORYINTEGRATION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterDirectoryIntegration(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERDIRECTORYINTEGRATION", settings)
}

// SetClusterDomainDefaults replaces the cluster-wide default Domain settings.
//
// Parameters:
//   - settings: a dictionary used to replace the default Domain settings.
//
// This method executes the SETCLUSTERDOMAINDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERDOMAINDEFAULTS", settings)
}

// SetClusterTrustedCerts replaces the cluster-wide list of Trusted Certificates.
//
// Parameters:
//   - certs: an array of Base64-encoded X.509 certificate data or PEM.
//
// This method executes the SETCLUSTERTRUSTEDCERTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterTrustedCerts(certs []string) error {
	if certs == nil {
		return fmt.Errorf("certs array is required")
	}
	return cli.setTrustedCerts("SETCLUSTERTRUSTEDCERTS", certs)
}

// SetDirectoryIntegration replaces the server-wide Directory Integration settings.
//
// Parameters:
//   - settings: a dictionary used to replace the Directory Integration settings.
//
// This method executes the SETDIRECTORYINTEGRATION CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDirectoryIntegration(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDIRECTORYINTEGRATION", settings)
}

// SetDomainDefaults replaces the server-wide default Domain settings.
//
// Parameters:
//   - settings: a dictionary used to replace the default Domain settings.
//
// This method executes the SETDOMAINDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDOMAINDEFAULTS", settings)
}

// SetServerAccountDefaults replaces the server-wide Default Account settings.
//
// Parameters:
//   - settings: a dictionary used to replace the Default Account settings.
//
// This method executes the SETSERVERACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetServerAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSERVERACCOUNTDEFAULTS", settings)
}

// SetServerAccountPrefs replaces the server-wide Default Account Preferences.
//
// Parameters:
//   - settings: a dictionary used to replace the Default Account Preferences.
//
// This method executes the SETSERVERACCOUNTPREFS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetServerAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSERVERACCOUNTPREFS", settings)
}

// SetServerTrustedCerts replaces the server-wide list of Trusted Certificates.
//
// Parameters:
//   - certs: an array of Base64-encoded X.509 certificate data or PEM.
//
// This method executes the SETSERVERTRUSTEDCERTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetServerTrustedCerts(certs []string) error {
	if certs == nil {
		return fmt.Errorf("certs array is required")
	}
	return cli.setTrustedCerts("SETSERVERTRUSTEDCERTS", certs)
}

// UpdateClusterAccountDefaults modifies specific cluster-wide Default Account settings.
//
// Parameters:
//   - settings: a dictionary containing settings to update.
//
// This method executes the UPDATECLUSTERACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateClusterAccountDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATECLUSTERACCOUNTDEFAULTS", settings)
}

// UpdateClusterAccountPrefs modifies specific cluster-wide Default Account Preferences.
//
// Parameters:
//   - settings: a dictionary containing preferences to update.
//
// This method executes the UPDATECLUSTERACCOUNTPREFS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateClusterAccountPrefs(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATECLUSTERACCOUNTPREFS", settings)
}

// UpdateClusterDomainDefaults modifies specific cluster-wide default Domain settings.
//
// Parameters:
//   - settings: a dictionary containing settings to update.
//
// This method executes the UPDATECLUSTERDOMAINDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateClusterDomainDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATECLUSTERDOMAINDEFAULTS", settings)
}

// UpdateDomainDefaults modifies specific server-wide default Domain settings.
//
// Parameters:
//   - settings: a dictionary containing settings to update.
//
// This method executes the UPDATEDOMAINDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateDomainDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEDOMAINDEFAULTS", settings)
}

// UpdateServerAccountDefaults modifies specific server-wide Default Account settings.
//
// Parameters:
//   - settings: a dictionary containing settings to update.
//
// This method executes the UPDATESERVERACCOUNTDEFAULTS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateServerAccountDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATESERVERACCOUNTDEFAULTS", settings)
}

// UpdateServerAccountPrefs modifies specific server-wide Default Account Preferences.
//
// Parameters:
//   - settings: a dictionary containing preferences to update.
//
// This method executes the UPDATESERVERACCOUNTPREFS CLI command.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateServerAccountPrefs(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATESERVERACCOUNTPREFS", settings)
}

func (cli *Cli) getTrustedCerts(cmd string) ([]string, error) {
	res, err := cli.getMapAny(cmd)
	if err != nil {
		return nil, err
	}
	val, ok := res["TrustedCertificates"]
	if !ok || val == nil {
		return []string{}, nil
	}
	raw, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("%s: expected []any for TrustedCertificates, got %T", cmd, val)
	}
	certs := make([]string, 0, len(raw))
	for i, item := range raw {
		b, ok := item.([]byte)
		if !ok {
			return nil, fmt.Errorf("%s: expected []byte at index %d, got %T", cmd, i, item)
		}
		certs = append(certs, base64.StdEncoding.EncodeToString(b))
	}
	return certs, nil
}

func (cli *Cli) listTelnums(cmd, filter string, limit int) (map[string]string, error) {
	args := make([]any, 0, 3)
	if filter != "" {
		args = append(args, "FILTER", filter)
	}
	if limit > 0 {
		args = append(args, limit)
	}
	res, err := cli.getMapAny(cmd, args...)
	if err != nil {
		return nil, err
	}
	delete(res, "")
	return toStringMap(cmd, res)
}

func (cli *Cli) setTrustedCerts(cmd string, certs []string) error {
	if certs == nil {
		return fmt.Errorf("%s: certs slice is nil", cmd)
	}
	replacer := strings.NewReplacer(
		"-----BEGIN CERTIFICATE-----", "",
		"-----END CERTIFICATE-----", "",
		"\n", "", "\r", "", " ", "",
	)
	prepared := make([]any, 0, len(certs))
	for _, c := range certs {
		clean := replacer.Replace(c)
		if clean != "" {
			prepared = append(prepared, DataBlock(clean))
		}
	}
	return cli.QueryNV(cmd, map[string]any{"TrustedCertificates": prepared})
}
