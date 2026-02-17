// # Server Settings
//
// The Server Settings section is the core of the library, providing comprehensive
// access to the CommuniGate Pro configuration engine. It allows for the
// management of both node-local and cluster-wide parameters, ensuring
// consistency across multi-server environments.
//
// Key capabilities include:
//   - Network Security: Management of various IP address lists including
//     Blacklisted, Client (Trusted), LAN, and Denied IPs for granular
//     access control.
//   - Routing & Rules: Integration with [RouterList], [MailRule], and
//     [SignalRule] at both Server and Cluster levels.
//   - Module Configuration: Dynamic retrieval and updates for specific
//     protocol modules (SMTP, IMAP, SIP, etc.) via [Cli.GetModule] and [Cli.UpdateModule].
//   - Content Filtering: Controlling prohibited content through the
//     [BannedLines] structure for headers and message bodies.
//   - Global Subsystems: Fine-tuning of DNR (Domain Name Resolver),
//     Logging, Queue Management, and Session settings.
//
// Use these methods to build automated configuration auditors,
// disaster recovery scripts, or centralized management consoles for
// large-scale CommuniGate Pro deployments.
package cgpcli

import (
	"fmt"
)

// BannedLines represents lists of prohibited strings in message headers and bodies.
type BannedLines struct {
	BodyLines   []string // A slice of strings to be banned in the message body.
	HeaderLines []string // A slice of strings to be banned in the message headers.
}

// GetBanned retrieves the server-wide Banned Message Lines settings.
//
// Returns:
//   - *[BannedLines]: a structure containing banned body and header lines.
//   - error: an error if the command fails.
func (cli *Cli) GetBanned() (*BannedLines, error) {
	return cli.getBannedGeneric("GETBANNED")
}

// GetBlacklistedIPs retrieves the server-wide Blacklisted IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetBlacklistedIPs() (IPList, error) {
	return cli.getIPList("GETBLACKLISTEDIPS")
}

// GetClientIPs retrieves the server-wide Client (trusted) IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClientIPs() (IPList, error) {
	return cli.getIPList("GETCLIENTIPS")
}

// GetClusterBanned retrieves the Cluster-wide Banned Message Lines settings.
//
// Returns:
//   - *BannedLines: a structure containing banned body and header lines.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterBanned() (*BannedLines, error) {
	return cli.getBannedGeneric("GETCLUSTERBANNED")
}

// GetClusterBlacklistedIPs retrieves the Cluster-wide Blacklisted IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterBlacklistedIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERBLACKLISTEDIPS")
}

// GetClusterClientIPs retrieves the Cluster-wide Client IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterClientIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERCLIENTIPS")
}

// GetClusterDebugIPs retrieves the Cluster-wide Debug IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterDebugIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERDEBUGIPS")
}

// GetClusterDeniedIPs retrieves the Cluster-wide Denied IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterDeniedIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERDENIEDIPS")
}

// GetClusterIntercept reads the Cluster-wide Lawful Intercept settings.
//
// Returns:
//   - map[string]any: a dictionary with the Intercept settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterIntercept() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERINTERCEPT")
}

// GetClusterLANIPs retrieves the Cluster-wide LAN IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterLANIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERLANIPS")
}

// GetClusterMailRules reads the Cluster-wide Automated Mail Processing Rules.
//
// Returns:
//   - [][MailRule]: a slice of Cluster Queue Rules.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterMailRules() ([]MailRule, error) {
	return cli.getMailRules("GETCLUSTERMAILRULES", "")
}

// GetClusterNatedIPs retrieves the Cluster-wide NATed IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterNatedIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERNATEDIPS")
}

// GetClusterNatSiteIPs retrieves the Cluster-wide NAT Site IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterNatSiteIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERNATSITEIPS")
}

// GetClusterNetwork retrieves the Cluster-wide Network settings.
//
// Returns:
//   - map[string]any: a dictionary with the network settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterNetwork() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERNETWORK")
}

// GetClusterRouterSettings reads the Cluster-wide Router settings.
//
// Returns:
//   - map[string]any: a dictionary with the Router settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterRouterSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERROUTERSETTINGS")
}

// GetClusterRouterTable reads the Cluster-wide Router Table.
//
// Returns:
//   - [RouterList]: a parsed list of routing rules.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterRouterTable() (RouterList, error) {
	return cli.getRouterTable("GETCLUSTERROUTERTABLE")
}

// GetClusterSettings retrieves the Cluster-wide settings.
//
// Returns:
//   - map[string]any: a dictionary with the Cluster settings.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERSETTINGS")
}

// GetClusterSignalRules reads the Cluster-wide Automated Signal Processing Rules.
//
// Returns:
//   - [][SignalRule]: a slice of Cluster Signal Rules.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterSignalRules() ([]SignalRule, error) {
	return cli.getSignalRules("GETCLUSTERSIGNALRULES", "")
}

// GetClusterWhiteHoleIPs retrieves the Cluster-wide WhiteHole IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetClusterWhiteHoleIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERWHITEHOLEIPS")
}

// GetDebugIPs retrieves the server-wide Debug IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetDebugIPs() (IPList, error) {
	return cli.getIPList("GETDEBUGIPS")
}

// GetDeniedIPs retrieves the server-wide Denied IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetDeniedIPs() (IPList, error) {
	return cli.getIPList("GETDENIEDIPS")
}

// GetDNRSettings retrieves the DNR (Domain Name Resolver) settings.
//
// Returns:
//   - map[string]any: a dictionary with the DNR settings.
//   - error: an error if the command fails.
func (cli *Cli) GetDNRSettings() (map[string]any, error) {
	return cli.getMapAny("GETDNRSETTINGS")
}

// GetIPState retrieves the type assigned to the specified IP address.
//
// Parameters:
//   - ip: specifies the IP address to check.
//   - temp: if true, the temporary Client IP Addresses set is checked.
//
// Returns:
//   - string: the IP address type.
//   - error: an error if the command fails.
func (cli *Cli) GetIPState(ip string, temp bool) (string, error) {
	const cmd = "GETIPSTATE"
	if temp {
		return cli.getString(cmd, ip, "TEMP")
	} else {
		return cli.getString(cmd, ip)
	}
}

// GetLANIPs retrieves the set of LAN IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetLANIPs() (IPList, error) {
	return cli.getIPList("GETLANIPS")
}

// GetLogSettings retrieves the Main Log settings.
//
// Returns:
//   - map[string]any: a dictionary with the Main Log settings.
//   - error: an error if the command fails.
func (cli *Cli) GetLogSettings() (map[string]any, error) {
	return cli.getMapAny("GETLOGSETTINGS")
}

// GetMediaServerSettings reads the Media Server component settings.
//
// Returns:
//   - map[string]any: a dictionary with the component settings.
//   - error: an error if the command fails.
func (cli *Cli) GetMediaServerSettings() (map[string]any, error) {
	return cli.getMapAny("GETMEDIASERVERSETTINGS")
}

// GetModule retrieves settings for a specific Server module.
//
// Parameters:
//   - name: specifies the name of a CommuniGate Pro Server module.
//
// Returns:
//   - map[string]any: a dictionary with the module settings.
//   - error: an error if the command fails.
func (cli *Cli) GetModule(name string) (map[string]any, error) {
	return cli.getMapAny("GETMODULE", name)
}

// GetNatedIPs retrieves the server-wide NATed IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetNatedIPs() (IPList, error) {
	return cli.getIPList("GETNATEDIPS")
}

// GetNatSiteIPs retrieves the server-wide NAT Site IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetNatSiteIPs() (IPList, error) {
	return cli.getIPList("GETNATSITEIPS")
}

// GetNetwork retrieves the Server Network settings.
//
// Returns:
//   - map[string]any: a dictionary with the server network settings.
//   - error: an error if the command fails.
func (cli *Cli) GetNetwork() (map[string]any, error) {
	return cli.getMapAny("GETNETWORK")
}

// GetQueueSettings retrieves the Queue settings.
//
// Returns:
//   - map[string]any: a dictionary with the Queue settings.
//   - error: an error if the command fails.
func (cli *Cli) GetQueueSettings() (map[string]any, error) {
	return cli.getMapAny("GETQUEUESETTINGS")
}

// GetRouterSettings reads the Router settings.
//
// Returns:
//   - map[string]any: a dictionary with the Router settings.
//   - error: an error if the command fails.
func (cli *Cli) GetRouterSettings() (map[string]any, error) {
	return cli.getMapAny("GETROUTERSETTINGS")
}

// GetRouterTable reads the Router Table.
//
// Returns:
//   - [RouterList]: a parsed list of routing rules.
//   - error: an error if the command fails.
func (cli *Cli) GetRouterTable() (RouterList, error) {
	return cli.getRouterTable("GETROUTERTABLE")
}

// GetServerIntercept reads the Lawful Intercept settings.
//
// Returns:
//   - map[string]any: a dictionary with the Intercept settings.
//   - error: an error if the command fails.
func (cli *Cli) GetServerIntercept() (map[string]any, error) {
	return cli.getMapAny("GETSERVERINTERCEPT")
}

// GetServerMailRules reads the Server-Wide Automated Mail Processing Rules.
//
// Returns:
//   - [][MailRule]: a slice of Server Queue Rules.
//   - error: an error if the command fails.
func (cli *Cli) GetServerMailRules() ([]MailRule, error) {
	return cli.getMailRules("GETSERVERMAILRULES", "")
}

// GetServerSettings reads the "other" Server settings.
//
// Returns:
//   - map[string]any: a dictionary with the Server settings.
//   - error: an error if the command fails.
func (cli *Cli) GetServerSettings() (map[string]any, error) {
	return cli.getMapAny("GETSERVERSETTINGS")
}

// GetServerSignalRules reads the Server-Wide Automated Signal Processing Rules.
//
// Returns:
//   - [][SignalRule]: a slice of Server Signal Rules.
//   - error: an error if the command fails.
func (cli *Cli) GetServerSignalRules() ([]SignalRule, error) {
	return cli.getSignalRules("GETSERVERSIGNALRULES", "")
}

// GetSessionSettings retrieves the user Sessions settings.
//
// Returns:
//   - map[string]any: a dictionary with the Sessions settings.
//   - error: an error if the command fails.
func (cli *Cli) GetSessionSettings() (map[string]any, error) {
	return cli.getMapAny("GETSESSIONSETTINGS")
}

// GetSignalSettings retrieves the Signal component settings.
//
// Returns:
//   - map[string]any: a dictionary with the component settings.
//   - error: an error if the command fails.
func (cli *Cli) GetSignalSettings() (map[string]any, error) {
	return cli.getMapAny("GETSIGNALSETTINGS")
}

// GetSMTPSendProfiles retrieves the Target Host Profiles of the SMTP module.
//
// Returns:
//   - map[string]any: a dictionary containing profile settings.
//   - error: an error if the command fails.
func (cli *Cli) GetSMTPSendProfiles() (map[string]any, error) {
	return cli.getMapAny("GETSMTPSENDPROFILES")
}

// GetWhiteHoleIPs retrieves the server-wide WhiteHole IP Addresses.
//
// Returns:
//   - [IPList]: a parsed list of IP addresses and ranges.
//   - error: an error if the command fails.
func (cli *Cli) GetWhiteHoleIPs() (IPList, error) {
	return cli.getIPList("GETWHITEHOLEIPS")
}

// ListModules lists all Server modules.
//
// Returns:
//   - []string: an array with all module names.
//   - error: an error if the command fails.
func (cli *Cli) ListModules() ([]string, error) {
	return cli.getSliceString("LISTMODULES")
}

// RefreshOSData makes the Server re-read IP and DNS data from the OS.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) RefreshOSData() error {
	return cli.QueryNV("REFRESHOSDATA")
}

// Route retrieves the routing for the specified address.
//
// Parameters:
//   - address: specifies the e-mail address to be processed.
//   - mode: (Optional) specifies the routing type (mail, access, or signal).
//     Defaults to access.
//
// Returns:
//   - []string: an array containing [module, host, address].
//   - error: an error if the command fails.
func (cli *Cli) Route(address string, mode string) ([]string, error) {
	if mode != "" {
		return cli.getSliceString("ROUTE", address, mode)
	} else {
		return cli.getSliceString("ROUTE", address)
	}
}

// SetBanned sets the server-wide Banned Message Lines settings.
//
// Parameters:
//   - settings: specifies new server Banned settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetBanned(settings *BannedLines) error {
	return cli.setBannedGeneric("SETBANNED", settings)
}

// SetBlacklistedIPs updates the server-wide Blacklisted IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetBlacklistedIPs(list IPList) error {
	return cli.QueryNV("SETBLACKLISTEDIPS", list)
}

// SetClientIPs updates the server-wide Client IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClientIPs(list IPList) error {
	return cli.QueryNV("SETCLIENTIPS", list)
}

// SetClusterBanned sets the Cluster-wide Banned Message Lines settings.
//
// Parameters:
//   - settings: specifies new cluster-wide Banned settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterBanned(settings *BannedLines) error {
	return cli.setBannedGeneric("SETCLUSTERBANNED", settings)
}

// SetClusterBlacklistedIPs updates the Cluster-wide Blacklisted IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterBlacklistedIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERBLACKLISTEDIPS", list)
}

// SetClusterClientIPs updates the Cluster-wide Client IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterClientIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERCLIENTIPS", list)
}

// SetClusterDebugIPs updates the Cluster-wide Debug IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterDebugIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERDEBUGIPS", list)
}

// SetClusterDeniedIPs updates the Cluster-wide Denied IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterDeniedIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERDENIEDIPS", list)
}

// SetClusterIntercept updates the Cluster-wide Lawful Intercept settings.
//
// Parameters:
//   - settings: a dictionary containing the new Intercept settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterIntercept(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERINTERCEPT", settings)
}

// SetClusterLANIPs updates the Cluster-wide LAN IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterLANIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERLANIPS", list)
}

// SetClusterMailRules updates the Cluster-wide Automated Mail Processing Rules.
//
// Parameters:
//   - rules: an array of new Cluster Queue Mail Rules.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterMailRules(rules []MailRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETCLUSTERMAILRULES", "", rules)
}

// SetClusterNatedIPs updates the Cluster-wide NATed IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterNatedIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERNATEDIPS", list)
}

// SetClusterNatSiteIPs updates the Cluster-wide NAT Site IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterNatSiteIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERNATSITEIPS", list)
}

// SetClusterNetwork updates the Cluster-wide Network settings.
//
// Parameters:
//   - settings: a dictionary containing the new network settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterNetwork(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERNETWORK", settings)
}

// SetClusterRouterSettings updates the Cluster-wide Router settings.
//
// Parameters:
//   - settings: a dictionary containing new Router settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterRouterSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERROUTERSETTINGS", settings)
}

// SetClusterRouterTable updates the Cluster-wide Router Table.
//
// Parameters:
//   - table: the new set of routing rules.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterRouterTable(table RouterList) error {
	return cli.QueryNV("SETCLUSTERROUTERTABLE", table)
}

// SetClusterSettings sets the Cluster-wide settings.
//
// Parameters:
//   - settings: a dictionary containing the new Cluster settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERSETTINGS", settings)
}

// SetClusterSignalRules updates the Cluster-wide Automated Signal Processing Rules.
//
// Parameters:
//   - rules: an array of new Cluster Signal Rules.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterSignalRules(rules []SignalRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETCLUSTERSIGNALRULES", "", rules)
}

// SetClusterWhiteHoleIPs updates the Cluster-wide WhiteHole IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetClusterWhiteHoleIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERWHITEHOLEIPS", list)
}

// SetDebugIPs updates the server-wide Debug IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDebugIPs(list IPList) error {
	return cli.QueryNV("SETDEBUGIPS", list)
}

// SetDeniedIPs updates the server-wide Denied IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDeniedIPs(list IPList) error {
	return cli.QueryNV("SETDENIEDIPS", list)
}

// SetDNRSettings updates the DNR (Domain Name Resolver) settings.
//
// Parameters:
//   - settings: a dictionary containing new DNR settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetDNRSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDNRSETTINGS", settings)
}

// SetLANIPs updates the set of LAN IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetLANIPs(list IPList) error {
	return cli.QueryNV("SETLANIPS", list)
}

// SetMediaServerSettings sets the Media Server component settings.
//
// Parameters:
//   - settings: a dictionary containing new Media Server settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetMediaServerSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETMEDIASERVERSETTINGS", settings)
}

// SetModule replaces the settings for a specific Server module.
//
// Parameters:
//   - module: specifies the name of a CommuniGate Pro Server module.
//   - settings: a dictionary containing the new module settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetModule(module string, settings map[string]any) error {
	if module == "" || settings == nil {
		return fmt.Errorf("module name and settings dictionary are required")
	}
	return cli.QueryNV("SETMODULE", module, settings)
}

// SetNatedIPs updates the server-wide NATed IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetNatedIPs(list IPList) error {
	return cli.QueryNV("SETNATEDIPS", list)
}

// SetNatSiteIPs updates the server-wide NAT Site IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetNatSiteIPs(list IPList) error {
	return cli.QueryNV("SETNATSITEIPS", list)
}

// SetNetwork updates the Server Network settings.
//
// Parameters:
//   - settings: a dictionary containing new network settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetNetwork(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETNETWORK", settings)
}

// SetQueueSettings sets the Queue settings.
//
// Parameters:
//   - settings: a dictionary containing new Queue settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetQueueSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETQUEUESETTINGS", settings)
}

// SetRouterSettings updates the Router settings.
//
// Parameters:
//   - settings: a dictionary containing new Router settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetRouterSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETROUTERSETTINGS", settings)
}

// SetRouterTable updates the Router Table.
//
// Parameters:
//   - table: the new set of routing rules.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetRouterTable(table RouterList) error {
	return cli.QueryNV("SETROUTERTABLE", table)
}

// SetServerIntercept updates the Lawful Intercept settings.
//
// Parameters:
//   - settings: a dictionary containing new Intercept settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetServerIntercept(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSERVERINTERCEPT", settings)
}

// SetServerMailRules updates the Server-Wide Automated Mail Processing Rules.
//
// Parameters:
//   - rules: an array of new Server Queue Rules.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetServerMailRules(rules []MailRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETSERVERMAILRULES", "", rules)
}

// SetServerSignalRules updates the Server-Wide Automated Signal Processing Rules.
//
// Parameters:
//   - rules: an array of new Server Signal Rules.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetServerSignalRules(rules []SignalRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETSERVERSIGNALRULES", "", rules)
}

// SetSessionSettings sets the user Sessions settings.
//
// Parameters:
//   - settings: a dictionary containing new Sessions settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetSessionSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSESSIONSETTINGS", settings)
}

// SetSignalSettings sets the Signal component settings.
//
// Parameters:
//   - settings: a dictionary containing new component settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetSignalSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSIGNALSETTINGS", settings)
}

// SetSMTPSendProfiles sets the Target Host Profiles of the SMTP module.
//
// Parameters:
//   - settings: a dictionary containing the new profile settings.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetSMTPSendProfiles(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSMTPSENDPROFILES", settings)
}

// SetWhiteHoleIPs updates the server-wide WhiteHole IP Addresses.
//
// Parameters:
//   - list: the new set of IP addresses and ranges.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) SetWhiteHoleIPs(list IPList) error {
	return cli.QueryNV("SETWHITEHOLEIPS", list)
}

// UpdateLogSettings modifies the Main Log settings.
//
// Parameters:
//   - settings: a dictionary containing settings to be updated.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateLogSettings(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATELOGSETTINGS", settings)
}

// UpdateModule modifies the settings for a specific Server module.
//
// Parameters:
//   - module: specifies the name of a CommuniGate Pro Server module.
//   - settings: a dictionary containing settings to be updated.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateModule(module string, settings map[string]any) error {
	if module == "" || len(settings) == 0 {
		return fmt.Errorf("module name and settings dictionary are required")
	}
	return cli.QueryNV("UPDATEMODULE", module, settings)
}

// UpdateServerSettings modifies the "other" Server settings.
//
// Parameters:
//   - settings: a dictionary containing settings to be updated.
//
// Returns:
//   - error: an error if the command fails.
func (cli *Cli) UpdateServerSettings(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATESERVERSETTINGS", settings)
}

func (cli *Cli) getBannedGeneric(cmd string) (*BannedLines, error) {
	m, err := cli.getMapAny(cmd)
	if err != nil {
		return nil, err
	}

	res := &BannedLines{}
	if v, ok := m["BannedBodyLines"]; ok {
		res.BodyLines, _ = toStringSlice(cmd, v)
	}
	if v, ok := m["BannedHeaderLines"]; ok {
		res.HeaderLines, _ = toStringSlice(cmd, v)
	}
	return res, nil
}

func (cli *Cli) getIPList(cmd string) (IPList, error) {
	s, err := cli.getString(cmd)
	if err != nil {
		return nil, err
	}
	return parseIPList(s), nil
}

func (cli *Cli) getRouterTable(cmd string) (RouterList, error) {
	raw, err := cli.getString(cmd)
	if err != nil {
		return nil, err
	}
	return parseRouterTable(raw), nil
}

func (cli *Cli) setBannedGeneric(cmd string, settings *BannedLines) error {
	if settings == nil {
		return fmt.Errorf("%s: BannedLines is nil", cmd)
	}
	data := map[string]any{
		"BannedBodyLines":   settings.BodyLines,
		"BannedHeaderLines": settings.HeaderLines,
	}
	return cli.QueryNV(cmd, data)
}
