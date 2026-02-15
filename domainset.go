// # Domain Set Administration
//
// The Domain Set section provides functionality for global configuration and
// high-level management of the CommuniGate Pro environment. It handles operations
// that affect the entire server or cluster, rather than specific individual domains.
//
// Key capabilities include:
//   - Domain Lifecycle: Creating and deleting domains (including [Directory Domains]),
//     renaming, and managing physical storage mount points.
//   - Cluster-Wide Defaults: Managing [Cluster Domain Defaults], [Account Defaults],
//     and [Preferences] that synchronize across all cluster nodes.
//   - Server-Wide Defaults: Configuring default settings and preferences used for
//     newly created objects at the single-server level.
//   - Security: Managing [Trusted Certificates] (CA lists) for both Cluster and Server levels.
//   - Integration: Configuring system-wide [Directory Integration] settings and
//     global telephony number (Telnums).
//
// Use these methods to maintain consistency across large-scale deployments and
// manage the infrastructure upon which individual domains reside.
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
//   - domain: This parameter specifies the Domain name to create.
//   - settings: An optional dictionary specifying the Domain settings.
//
// Returns:
//   - error: An error if the domain name is empty or the command fails.
func (cli *Cli) CreateDirectoryDomain(domain string, settings map[string]any) (err error) {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "CREATEDIRECTORYDOMAIN"
	if settings != nil {
		err = cli.QueryNV(cmd, domain, settings)
	} else {
		err = cli.QueryNV(cmd, domain)
	}
	return
}

// CreateDomain creates a new secondary Domain.
//
// Parameters:
//   - domain: This parameter specifies the Domain name to create.
//   - shared: If true, creates a Cluster-wide Domain (Dynamic Cluster only).
//   - storage: Optional "storage mount Point" directory name (without .mnt suffix).
//   - settings: An optional dictionary specifying the initial Domain settings.
//
// Returns:
//   - error: An error if the domain name is empty or the server fails to create the domain.
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
//   - storage: This parameter specifies the "storage mount Point" name.
//   - shared: If true, creates a mount point for Cluster Domains in a Dynamic Cluster.
//
// Returns:
//   - error: An error if the storage path is empty or the command fails.
func (cli *Cli) CreateDomainStorage(storage string, shared bool) (err error) {
	if storage == "" {
		return fmt.Errorf("storage path is required")
	}

	const cmd = "CREATEDOMAINSTORAGE"
	if shared {
		err = cli.QueryNV(cmd, "SHARED", "PATH", storage)
	} else {
		err = cli.QueryNV(cmd, "PATH", storage)
	}
	return
}

// DeleteDomain removes an existing Domain.
//
// Parameters:
//   - domain: This parameter specifies the name of the Domain to be removed.
//   - force: If true, removes the Domain even if it contains accounts or other objects.
//
// Returns:
//   - error: An error if the domain name is empty or the server fails to remove the domain.
func (cli *Cli) DeleteDomain(domain string, force bool) (err error) {
	if domain == "" {
		return fmt.Errorf("domain name is required")
	}

	const cmd = "DELETEDOMAIN"
	if force {
		err = cli.QueryNV(cmd, domain, "FORCE")
	} else {
		err = cli.QueryNV(cmd, domain)
	}
	return
}

// GetClusterAccountDefaults retrieves the cluster-wide Default Account settings.
//
// Returns:
//   - map[string]any: A dictionary with the cluster default Account settings.
//   - error: An error if the command fails or is not executed in a Dynamic Cluster.
func (cli *Cli) GetClusterAccountDefaults() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERACCOUNTDEFAULTS")
}

// GetClusterAccountPrefs retrieves the cluster-wide Default Account Preferences.
//
// Returns:
//   - map[string]any: A dictionary with the cluster default Preferences.
//   - error: An error if the command fails or is not executed in a Dynamic Cluster.
func (cli *Cli) GetClusterAccountPrefs() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERACCOUNTPREFS")
}

// GetClusterDirectoryIntegration retrieves the cluster-wide Directory Integration settings.
//
// Returns:
//   - map[string]any: A dictionary with the cluster Directory Integration settings.
//   - error: An error if the command fails or is not executed in a Dynamic Cluster.
func (cli *Cli) GetClusterDirectoryIntegration() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDIRECTORYINTEGRATION")
}

// GetClusterDomainDefaults retrieves the cluster-wide default Domain Settings.
//
// Returns:
//   - map[string]any: A dictionary with the cluster-wide default Domain Settings.
//   - error: An error if the command fails or is not executed in a Dynamic Cluster.
func (cli *Cli) GetClusterDomainDefaults() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDOMAINDEFAULTS")
}

// GetClusterTrustedCerts retrieves the cluster-wide set of Trusted Certificates.
//
// Returns:
//   - []string: A slice of Base64-encoded X.509 certificate data.
//   - error: An error if the command fails.
func (cli *Cli) GetClusterTrustedCerts() ([]string, error) {
	return cli.getTrustedCerts("GETCLUSTERTRUSTEDCERTS")
}

// GetDirectoryIntegration retrieves the server-wide Directory Integration settings.
//
// Returns:
//   - map[string]any: A dictionary with the Directory Integration settings.
//   - error: An error if the command fails.
func (cli *Cli) GetDirectoryIntegration() (map[string]any, error) {
	return cli.getMapAny("GETDIRECTORYINTEGRATION")
}

// GetDomainDefaults retrieves the server-wide default Domain Settings.
//
// Returns:
//   - map[string]any: A dictionary with the default Domain Settings.
//   - error: An error if the command fails.
func (cli *Cli) GetDomainDefaults() (map[string]any, error) {
	return cli.getMapAny("GETDOMAINDEFAULTS")
}

// GetServerAccountDefaults retrieves the server-wide Default Account settings.
//
// Returns:
//   - map[string]any: A dictionary with the default Account settings.
//   - error: An error if the command fails.
func (cli *Cli) GetServerAccountDefaults() (map[string]any, error) {
	return cli.getMapAny("GETSERVERACCOUNTDEFAULTS")
}

// GetServerAccountPrefs retrieves the server-wide Default Account Preferences.
//
// Returns:
//   - map[string]any: A dictionary with the default Preferences.
//   - error: An error if the command fails.
func (cli *Cli) GetServerAccountPrefs() (map[string]any, error) {
	return cli.getMapAny("GETSERVERACCOUNTPREFS")
}

// GetServerTrustedCerts retrieves the server-wide set of Trusted Certificates.
//
// Returns:
//   - []string: An array of Base64-encoded X.509 certificate data.
//   - error: An error if the command fails.
func (cli *Cli) GetServerTrustedCerts() ([]string, error) {
	return cli.getTrustedCerts("GETSERVERTRUSTEDCERTS")
}

// ListClusterTelnums reads Telnum numbers created in shared Cluster Domains.
//
// Parameters:
//   - filter: (Optional) substring filter for telnum numbers.
//   - limit: Maximum number of Telnum numbers to return.
//
// Returns:
//   - map[string]string: A dictionary mapping Telnums to Account names.
//   - error: An error if the command fails.
func (cli *Cli) ListClusterTelnums(filter string, limit int) (map[string]string, error) {
	return cli.listTelnums("LISTCLUSTERTELNUMS", filter, limit)
}

// ListDomainStorage retrieves the list of "storage mount points" for Domains.
//
// Parameters:
//   - shared: If true, lists mount points for Cluster Domains.
//
// Returns:
//   - []string: An array containing the "storage mount points" names.
//   - error: An error if the command fails.
func (cli *Cli) ListDomainStorage(shared bool) (res []string, err error) {
	const cmd = "LISTDOMAINSTORAGE"
	if shared {
		res, err = cli.getSliceString(cmd, "SHARED")
	} else {
		res, err = cli.getSliceString(cmd)
	}
	return
}

// ListDomains retrieves the names of all server domains.
//
// Returns:
//   - []string: An array with the names of all server domains.
//   - error: An error if the command fails.
func (cli *Cli) ListDomains() ([]string, error) {
	return cli.getSliceString("LISTDOMAINS")
}

// ListServerTelnums reads Telnum numbers created in all non-clustered Domains.
//
// Parameters:
//   - filter: (Optional) substring filter for telnum numbers.
//   - limit: Maximum number of Telnum numbers to return.
//
// Returns:
//   - map[string]string: A dictionary mapping Telnums to Account names.
//   - error: An error if the command fails.
func (cli *Cli) ListServerTelnums(filter string, limit int) (map[string]string, error) {
	return cli.listTelnums("LISTSERVERTELNUMS", filter, limit)
}

// MainDomainName retrieves the name of the Main Domain.
//
// Returns:
//   - string: The Main Domain name.
//   - error: An error if the command fails.
func (cli *Cli) MainDomainName() (string, error) {
	return cli.getString("MAINDOMAINNAME")
}

// ReloadDirectoryDomains forces the server to scan the Domains Directory subtree.
// This allows finding additional Directory-based Domains created bypassing the CLI.
//
// Returns:
//   - error: An error if the operation is not enabled or fails.
func (cli *Cli) ReloadDirectoryDomains() error {
	return cli.QueryNV("RELOADDIRECTORYDOMAINS")
}

// RenameDomain changes the name of an existing secondary Domain.
//
// Parameters:
//   - oldName: The name of an existing secondary Domain.
//   - newName: The new name for the Domain.
//   - storage: Optional new "storage mount Point" directory for the Domain data.
//
// Returns:
//   - error: An error if names are missing or the server fails to rename the domain.
func (cli *Cli) RenameDomain(oldName, newName string, storage string) (err error) {
	if oldName == "" || newName == "" {
		return fmt.Errorf("old and new domain names are required")
	}

	const cmd = "RENAMEDOMAIN"
	if storage != "" {
		err = cli.QueryNV(cmd, oldName, "INTO", newName, "PATH", storage)
	} else {
		err = cli.QueryNV(cmd, oldName, "INTO", newName)
	}
	return
}

// SetClusterAccountDefaults replaces the cluster-wide Default Account settings.
//
// Parameters:
//   - settings: A dictionary used to replace the Default Account settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetClusterAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERACCOUNTDEFAULTS", settings)
}

// SetClusterAccountPrefs replaces the cluster-wide Default Account Preferences.
//
// Parameters:
//   - settings: A dictionary used to replace the Default Account Preferences.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetClusterAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERACCOUNTPREFS", settings)
}

// SetClusterDirectoryIntegration replaces the cluster-wide Directory Integration settings.
//
// Parameters:
//   - settings: A dictionary used to replace the Directory Integration settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetClusterDirectoryIntegration(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERDIRECTORYINTEGRATION", settings)
}

// SetClusterDomainDefaults replaces the cluster-wide default Domain settings.
//
// Parameters:
//   - settings: A dictionary used to replace the default Domain settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetClusterDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERDOMAINDEFAULTS", settings)
}

// SetClusterTrustedCerts replaces the cluster-wide list of Trusted Certificates.
//
// Parameters:
//   - certs: An array of Base64-encoded X.509 certificate data or PEM.
//
// Returns:
//   - error: An error if the input is nil or the command fails.
func (cli *Cli) SetClusterTrustedCerts(certs []string) error {
	if certs == nil {
		return fmt.Errorf("certs array is required")
	}
	return cli.setTrustedCerts("SETCLUSTERTRUSTEDCERTS", certs)
}

// SetDirectoryIntegration replaces the server-wide Directory Integration settings.
//
// Parameters:
//   - settings: A dictionary used to replace the Directory Integration settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetDirectoryIntegration(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDIRECTORYINTEGRATION", settings)
}

// SetDomainDefaults replaces the server-wide default Domain settings.
//
// Parameters:
//   - settings: A dictionary used to replace the default Domain settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetDomainDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDOMAINDEFAULTS", settings)
}

// SetServerAccountDefaults replaces the server-wide Default Account settings.
//
// Parameters:
//   - settings: A dictionary used to replace the Default Account settings.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetServerAccountDefaults(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSERVERACCOUNTDEFAULTS", settings)
}

// SetServerAccountPrefs replaces the server-wide Default Account Preferences.
//
// Parameters:
//   - settings: A dictionary used to replace the Default Account Preferences.
//
// Returns:
//   - error: An error if settings are nil or the command fails.
func (cli *Cli) SetServerAccountPrefs(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSERVERACCOUNTPREFS", settings)
}

// SetServerTrustedCerts replaces the server-wide list of Trusted Certificates.
//
// Parameters:
//   - certs: An array of Base64-encoded X.509 certificate data or PEM.
//
// Returns:
//   - error: An error if the input is nil or the command fails.
func (cli *Cli) SetServerTrustedCerts(certs []string) error {
	if certs == nil {
		return fmt.Errorf("certs array is required")
	}
	return cli.setTrustedCerts("SETSERVERTRUSTEDCERTS", certs)
}

// UpdateClusterAccountDefaults modifies specific cluster-wide Default Account settings.
//
// Parameters:
//   - settings: A dictionary containing settings to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateClusterAccountDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATECLUSTERACCOUNTDEFAULTS", settings)
}

// UpdateClusterAccountPrefs modifies specific cluster-wide Default Account Preferences.
//
// Parameters:
//   - settings: A dictionary containing preferences to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateClusterAccountPrefs(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATECLUSTERACCOUNTPREFS", settings)
}

// UpdateClusterDomainDefaults modifies specific cluster-wide default Domain settings.
//
// Parameters:
//   - settings: A dictionary containing settings to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateClusterDomainDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATECLUSTERDOMAINDEFAULTS", settings)
}

// UpdateDomainDefaults modifies specific server-wide default Domain settings.
//
// Parameters:
//   - settings: A dictionary containing settings to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateDomainDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATEDOMAINDEFAULTS", settings)
}

// UpdateServerAccountDefaults modifies specific server-wide Default Account settings.
//
// Parameters:
//   - settings: A dictionary containing settings to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
func (cli *Cli) UpdateServerAccountDefaults(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATESERVERACCOUNTDEFAULTS", settings)
}

// UpdateServerAccountPrefs modifies specific server-wide Default Account Preferences.
//
// Parameters:
//   - settings: A dictionary containing preferences to update.
//
// Returns:
//   - error: An error if settings are empty or the command fails.
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
