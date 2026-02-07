// # Section: Server Settings
package cgpcli

import (
	"fmt"
)

// BannedLines представляет списки запрещенных строк в заголовках и теле сообщений.
type BannedLines struct {
	BodyLines   []string
	HeaderLines []string
}

// GetBanned возвращает серверные списки запрещенных строк.
func (cli *Cli) GetBanned() (*BannedLines, error) {
	return cli.getBannedGeneric("GETBANNED")
}

func (cli *Cli) GetBlacklistedIPs() (IPList, error) {
	return cli.getIPList("GETBLACKLISTEDIPS")
}

func (cli *Cli) GetClientIPs() (IPList, error) {
	return cli.getIPList("GETCLIENTIPS")
}

// GetClusterBanned возвращает общекластерные списки запрещенных строк.
func (cli *Cli) GetClusterBanned() (*BannedLines, error) {
	return cli.getBannedGeneric("GETCLUSTERBANNED")
}

func (cli *Cli) GetClusterBlacklistedIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERBLACKLISTEDIPS")
}

func (cli *Cli) GetClusterClientIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERCLIENTIPS")
}

func (cli *Cli) GetClusterDebugIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERDEBUGIPS")
}

func (cli *Cli) GetClusterDeniedIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERDENIEDIPS")
}

func (cli *Cli) GetClusterIntercept() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERINTERCEPT")
}

func (cli *Cli) GetClusterLANIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERLANIPS")
}

func (cli *Cli) GetClusterMailRules() ([]MailRule, error) {
	return cli.getMailRules("GETCLUSTERMAILRULES", "")
}

func (cli *Cli) GetClusterNatedIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERNATEDIPS")
}

func (cli *Cli) GetClusterNatSiteIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERNATSITEIPS")
}

func (cli *Cli) GetClusterNetwork() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERNETWORK")
}

func (cli *Cli) GetClusterRouterSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERROUTERSETTINGS")
}

func (cli *Cli) GetClusterRouterTable() (RouterList, error) {
	return cli.getRouterTable("GETCLUSTERROUTERTABLE")
}

func (cli *Cli) GetClusterSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERSETTINGS")
}

func (cli *Cli) GetClusterSignalRules() ([]SignalRule, error) {
	return cli.getSignalRules("GETCLUSTERSIGNALRULES", "")
}

func (cli *Cli) GetClusterWhiteHoleIPs() (IPList, error) {
	return cli.getIPList("GETCLUSTERWHITEHOLEIPS")
}

func (cli *Cli) GetDebugIPs() (IPList, error) {
	return cli.getIPList("GETDEBUGIPS")
}

func (cli *Cli) GetDeniedIPs() (IPList, error) {
	return cli.getIPList("GETDENIEDIPS")
}

func (cli *Cli) GetDNRSettings() (map[string]any, error) {
	return cli.getMapAny("GETDNRSETTINGS")
}

func (cli *Cli) GetIPState(ip string, temp bool) (s string, err error) {
	const cmd = "GETIPSTATE"
	if temp {
		s, err = cli.getString(cmd, ip, "TEMP")
	} else {
		s, err = cli.getString(cmd, ip)
	}
	return
}

func (cli *Cli) GetLANIPs() (IPList, error) {
	return cli.getIPList("GETLANIPS")
}

func (cli *Cli) GetLogSettings() (map[string]any, error) {
	return cli.getMapAny("GETLOGSETTINGS")
}

func (cli *Cli) GetMediaServerSettings() (map[string]any, error) {
	return cli.getMapAny("GETMEDIASERVERSETTINGS")
}

func (cli *Cli) GetModule(name string) (map[string]any, error) {
	return cli.getMapAny("GETMODULE", name)
}

func (cli *Cli) GetNatedIPs() (IPList, error) {
	return cli.getIPList("GETNATEDIPS")
}

func (cli *Cli) GetNatSiteIPs() (IPList, error) {
	return cli.getIPList("GETNATSITEIPS")
}

func (cli *Cli) GetNetwork() (map[string]any, error) {
	return cli.getMapAny("GETNETWORK")
}

func (cli *Cli) GetQueueSettings() (map[string]any, error) {
	return cli.getMapAny("GETQUEUESETTINGS")
}

func (cli *Cli) GetRouterSettings() (map[string]any, error) {
	return cli.getMapAny("GETROUTERSETTINGS")
}

func (cli *Cli) GetRouterTable() (RouterList, error) {
	return cli.getRouterTable("GETROUTERTABLE")
}

func (cli *Cli) GetServerIntercept() (map[string]any, error) {
	return cli.getMapAny("GETSERVERINTERCEPT")
}

func (cli *Cli) GetServerMailRules() ([]MailRule, error) {
	return cli.getMailRules("GETSERVERMAILRULES", "")
}

func (cli *Cli) GetServerSettings() (map[string]any, error) {
	return cli.getMapAny("GETSERVERSETTINGS")
}

func (cli *Cli) GetServerSignalRules() ([]SignalRule, error) {
	return cli.getSignalRules("GETSERVERSIGNALRULES", "")
}

func (cli *Cli) GetSessionSettings() (map[string]any, error) {
	return cli.getMapAny("GETSESSIONSETTINGS")
}

func (cli *Cli) GetSignalSettings() (map[string]any, error) {
	return cli.getMapAny("GETSIGNALSETTINGS")
}

func (cli *Cli) GetSMTPSendProfiles() (map[string]any, error) {
	return cli.getMapAny("GETSMTPSENDPROFILES")
}

func (cli *Cli) GetWhiteHoleIPs() (IPList, error) {
	return cli.getIPList("GETWHITEHOLEIPS")
}

func (cli *Cli) ListModules() ([]string, error) {
	return cli.getSliceString("LISTMODULES")
}

func (cli *Cli) RefreshOSData() error {
	return cli.QueryNV("REFRESHOSDATA")
}

func (cli *Cli) Route(address string, mode string) (res []string, err error) {
	if mode != "" {
		res, err = cli.getSliceString("ROUTE", address, mode)
	} else {
		res, err = cli.getSliceString("ROUTE", address)
	}
	return
}

// SetBanned устанавливает серверные списки запрещенных строк.
func (cli *Cli) SetBanned(lines *BannedLines) error {
	return cli.setBannedGeneric("SETBANNED", lines)
}

func (cli *Cli) SetBlacklistedIPs(list IPList) error {
	return cli.QueryNV("SETBLACKLISTEDIPS", list)
}

func (cli *Cli) SetClientIPs(list IPList) error {
	return cli.QueryNV("SETCLIENTIPS", list)
}

// SetClusterBanned устанавливает общекластерные списки запрещенных строк.
func (cli *Cli) SetClusterBanned(lines *BannedLines) error {
	return cli.setBannedGeneric("SETCLUSTERBANNED", lines)
}

func (cli *Cli) SetClusterBlacklistedIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERBLACKLISTEDIPS", list)
}

func (cli *Cli) SetClusterClientIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERCLIENTIPS", list)
}

func (cli *Cli) SetClusterDebugIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERDEBUGIPS", list)
}

func (cli *Cli) SetClusterDeniedIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERDENIEDIPS", list)
}

func (cli *Cli) SetClusterIntercept(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERINTERCEPT", settings)
}

func (cli *Cli) SetClusterLANIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERLANIPS", list)
}

func (cli *Cli) SetClusterMailRules(rules []MailRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETCLUSTERMAILRULES", "", rules)
}

func (cli *Cli) SetClusterNatedIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERNATEDIPS", list)
}

func (cli *Cli) SetClusterNatSiteIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERNATSITEIPS", list)
}

func (cli *Cli) SetClusterNetwork(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERNETWORK", settings)
}

func (cli *Cli) SetClusterRouterSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERROUTERSETTINGS", settings)
}

func (cli *Cli) SetClusterRouterTable(table RouterList) error {
	return cli.QueryNV("SETCLUSTERROUTERTABLE", table)
}

func (cli *Cli) SetClusterSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETCLUSTERSETTINGS", settings)
}

func (cli *Cli) SetClusterSignalRules(rules []SignalRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETCLUSTERSIGNALRULES", "", rules)
}

func (cli *Cli) SetClusterWhiteHoleIPs(list IPList) error {
	return cli.QueryNV("SETCLUSTERWHITEHOLEIPS", list)
}

func (cli *Cli) SetDebugIPs(list IPList) error {
	return cli.QueryNV("SETDEBUGIPS", list)
}

func (cli *Cli) SetDeniedIPs(list IPList) error {
	return cli.QueryNV("SETDENIEDIPS", list)
}

func (cli *Cli) SetDNRSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETDNRSETTINGS", settings)
}

func (cli *Cli) SetLANIPs(list IPList) error {
	return cli.QueryNV("SETLANIPS", list)
}

func (cli *Cli) SetMediaServerSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETMEDIASERVERSETTINGS", settings)
}

func (cli *Cli) SetModule(module string, settings map[string]any) error {
	if module == "" || settings == nil {
		return fmt.Errorf("module name and settings dictionary are required")
	}
	return cli.QueryNV("SETMODULE", module, settings)
}

func (cli *Cli) SetNatedIPs(list IPList) error {
	return cli.QueryNV("SETNATEDIPS", list)
}

func (cli *Cli) SetNatSiteIPs(list IPList) error {
	return cli.QueryNV("SETNATSITEIPS", list)
}

func (cli *Cli) SetNetwork(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETNETWORK", settings)
}

func (cli *Cli) SetQueueSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETQUEUESETTINGS", settings)
}

func (cli *Cli) SetRouterSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETROUTERSETTINGS", settings)
}

func (cli *Cli) SetRouterTable(table RouterList) error {
	return cli.QueryNV("SETROUTERTABLE", table)
}

func (cli *Cli) SetServerIntercept(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSERVERINTERCEPT", settings)
}

func (cli *Cli) SetServerMailRules(rules []MailRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETSERVERMAILRULES", "", rules)
}

func (cli *Cli) SetServerSignalRules(rules []SignalRule) error {
	if rules == nil {
		return fmt.Errorf("rules array is required")
	}
	return setRulesGeneric(cli, "SETSERVERSIGNALRULES", "", rules)
}

func (cli *Cli) SetSessionSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSESSIONSETTINGS", settings)
}

func (cli *Cli) SetSignalSettings(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSIGNALSETTINGS", settings)
}

func (cli *Cli) SetSMTPSendProfiles(settings map[string]any) error {
	if settings == nil {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("SETSMTPSENDPROFILES", settings)
}

func (cli *Cli) SetWhiteHoleIPs(list IPList) error {
	return cli.QueryNV("SETWHITEHOLEIPS", list)
}

func (cli *Cli) UpdateLogSettings(settings map[string]any) error {
	if len(settings) == 0 {
		return fmt.Errorf("settings dictionary is required")
	}
	return cli.QueryNV("UPDATELOGSETTINGS", settings)
}

func (cli *Cli) UpdateModule(module string, settings map[string]any) error {
	if module == "" || len(settings) == 0 {
		return fmt.Errorf("module name and settings dictionary are required")
	}
	return cli.QueryNV("UPDATEMODULE", module, settings)
}

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

func (cli *Cli) setBannedGeneric(cmd string, lines *BannedLines) error {
	if lines == nil {
		return fmt.Errorf("%s: BannedLines is nil", cmd)
	}
	data := map[string]any{
		"BannedBodyLines":   lines.BodyLines,
		"BannedHeaderLines": lines.HeaderLines,
	}
	return cli.QueryNV(cmd, data)
}
