// # Section: Server Settings
package cgpcli

import (
	"fmt"
	"strings"
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

func (cli *Cli) GetBlacklistedIPs() (IPList, error) { return cli.getIPListTyped("GETBLACKLISTEDIPS") }
func (cli *Cli) GetClientIPs() (IPList, error)      { return cli.getIPListTyped("GETCLIENTIPS") }

// GetClusterBanned возвращает общекластерные списки запрещенных строк.
func (cli *Cli) GetClusterBanned() (*BannedLines, error) {
	return cli.getBannedGeneric("GETCLUSTERBANNED")
}

func (cli *Cli) GetClusterBlacklistedIPs() (IPList, error) {
	return cli.getIPListTyped("GETCLUSTERBLACKLISTEDIPS")
}
func (cli *Cli) GetClusterClientIPs() (IPList, error) {
	return cli.getIPListTyped("GETCLUSTERCLIENTIPS")
}
func (cli *Cli) GetClusterDebugIPs() (IPList, error) { return cli.getIPListTyped("GETCLUSTERDEBUGIPS") }
func (cli *Cli) GetClusterDeniedIPs() (IPList, error) {
	return cli.getIPListTyped("GETCLUSTERDENIEDIPS")
}

func (cli *Cli) GetClusterDNRSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERDNRSETTINGS", "")
}

func (cli *Cli) GetClusterIntercept() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERINTERCEPT", "")
}

func (cli *Cli) GetClusterLANIPs() (IPList, error) { return cli.getIPListTyped("GETCLUSTERLANIPS") }

func (cli *Cli) GetClusterMailRules() ([]MailRule, error) {
	return cli.getMailRules("GETCLUSTERMAILRULES", "")
}

func (cli *Cli) GetClusterNatedIPs() (IPList, error) { return cli.getIPListTyped("GETCLUSTERNATEDIPS") }
func (cli *Cli) GetClusterNatSiteIPs() (IPList, error) {
	return cli.getIPListTyped("GETCLUSTERNATSITEIPS")
}

func (cli *Cli) GetClusterNetwork() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERNETWORK", "")
}

func (cli *Cli) GetClusterRouterSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERROUTERSETTINGS", "")
}

func (cli *Cli) GetClusterRouterTable() (string, error) {
	return cli.getString("GETCLUSTERROUTERTABLE")
}

func (cli *Cli) GetClusterSettings() (map[string]any, error) {
	return cli.getMapAny("GETCLUSTERSETTINGS", "")
}

func (cli *Cli) GetClusterSignalRules() ([]SignalRule, error) {
	return cli.getSignalRules("GETCLUSTERSIGNALRULES", "")
}

func (cli *Cli) GetClusterWhiteHoleIPs() (IPList, error) {
	return cli.getIPListTyped("GETCLUSTERWHITEHOLEIPS")
}

func (cli *Cli) GetDebugIPs() (IPList, error)  { return cli.getIPListTyped("GETDEBUGIPS") }
func (cli *Cli) GetDeniedIPs() (IPList, error) { return cli.getIPListTyped("GETDENIEDIPS") }

func (cli *Cli) GetDNRSettings() (map[string]any, error) {
	return cli.getMapAny("GETDNRSETTINGS", "")
}

func (cli *Cli) GetIPState(ip string, temp bool) (string, error) {
	if temp {
		return cli.getString("GETIPSTATE", ip, "TEMP")
	}
	return cli.getString("GETIPSTATE", ip)
}

func (cli *Cli) GetLANIPs() (IPList, error) { return cli.getIPListTyped("GETLANIPS") }

func (cli *Cli) GetLogSettings() (map[string]any, error) { return cli.getMapAny("GETLOGSETTINGS", "") }

func (cli *Cli) GetMediaServerSettings() (map[string]any, error) {
	return cli.getMapAny("GETMEDIASERVERSETTINGS", "")
}

func (cli *Cli) GetModule(name string) (map[string]any, error) {
	return cli.getMapAny("GETMODULE", name)
}

func (cli *Cli) GetNatedIPs() (IPList, error)        { return cli.getIPListTyped("GETNATEDIPS") }
func (cli *Cli) GetNatSiteIPs() (IPList, error)      { return cli.getIPListTyped("GETNATSITEIPS") }
func (cli *Cli) GetNetwork() (map[string]any, error) { return cli.getMapAny("GETNETWORK", "") }

func (cli *Cli) GetQueueSettings() (map[string]any, error) {
	return cli.getMapAny("GETQUEUESETTINGS", "")
}

func (cli *Cli) GetRouterSettings() (map[string]any, error) {
	return cli.getMapAny("GETROUTERSETTINGS", "")
}

func (cli *Cli) GetRouterTable() (RouterList, error) {
	raw, err := cli.getString("GETROUTERTABLE")
	if err != nil {
		return nil, err
	}

	parts := strings.Split(raw, `\e`)

	anyParts := make([]any, len(parts))
	for i, v := range parts {
		anyParts[i] = v
	}

	return parseRouterTable(anyParts), nil
}

func (cli *Cli) GetServerIntercept() (map[string]any, error) {
	return cli.getMapAny("GETSERVERINTERCEPT", "")
}

func (cli *Cli) GetServerMailRules() ([]MailRule, error) {
	return cli.getMailRules("GETSERVERMAILRULES", "")
}

func (cli *Cli) GetServerSettings() (map[string]any, error) {
	return cli.getMapAny("GETSERVERSETTINGS", "")
}

func (cli *Cli) GetServerSignalRules() ([]SignalRule, error) {
	return cli.getSignalRules("GETSERVERSIGNALRULES", "")
}

func (cli *Cli) GetSessionSettings() (map[string]any, error) {
	return cli.getMapAny("GETSESSIONSETTINGS", "")
}

func (cli *Cli) GetSignalSettings() (map[string]any, error) {
	return cli.getMapAny("GETSIGNALSETTINGS", "")
}

func (cli *Cli) GetSMTPSendProfiles() (map[string]any, error) {
	return cli.getMapAny("GETSMTPSENDPROFILES", "")
}

func (cli *Cli) GetWhiteHoleIPs() (IPList, error) { return cli.getIPListTyped("GETWHITEHOLEIPS") }
func (cli *Cli) ListModules() ([]string, error)   { return cli.getSliceString("LISTMODULES", "") }

func (cli *Cli) RefreshOSData() error {
	return cli.QueryNV("REFRESHOSDATA")
}

func (cli *Cli) Route(address string, mode string) ([]string, error) {
	var res any
	var err error
	if mode != "" {
		res, err = cli.Query("ROUTE", address, Atom(mode))
	} else {
		res, err = cli.Query("ROUTE", address)
	}
	if err != nil {
		return nil, err
	}
	return toStringSlice("ROUTE", res)
}

// SetBanned устанавливает серверные списки запрещенных строк.
func (cli *Cli) SetBanned(lines *BannedLines) error {
	return cli.setBannedGeneric("SETBANNED", lines)
}

func (cli *Cli) SetBlacklistedIPs(l IPList) error { return cli.setIPListTyped("SETBLACKLISTEDIPS", l) }
func (cli *Cli) SetClientIPs(l IPList) error      { return cli.setIPListTyped("SETCLIENTIPS", l) }

// SetClusterBanned устанавливает общекластерные списки запрещенных строк.
func (cli *Cli) SetClusterBanned(lines *BannedLines) error {
	return cli.setBannedGeneric("SETCLUSTERBANNED", lines)
}

func (cli *Cli) SetClusterBlacklistedIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERBLACKLISTEDIPS", l)
}
func (cli *Cli) SetClusterClientIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERCLIENTIPS", l)
}
func (cli *Cli) SetClusterDebugIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERDEBUGIPS", l)
}
func (cli *Cli) SetClusterDeniedIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERDENIEDIPS", l)
}

func (cli *Cli) SetClusterDNRSettings(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERDNRSETTINGS", "", s)
}

func (cli *Cli) SetClusterIntercept(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERINTERCEPT", "", s)
}

func (cli *Cli) SetClusterLANIPs(l IPList) error { return cli.setIPListTyped("SETCLUSTERLANIPS", l) }

func (cli *Cli) SetClusterMailRules(r []MailRule) error {
	return setRulesGeneric(cli, "SETCLUSTERMAILRULES", "", r)
}

func (cli *Cli) SetClusterNatedIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERNATEDIPS", l)
}
func (cli *Cli) SetClusterNatSiteIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERNATSITEIPS", l)
}

func (cli *Cli) SetClusterNetwork(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERNETWORK", "", s)
}

func (cli *Cli) SetClusterRouterSettings(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERROUTERSETTINGS", "", s)
}

func (cli *Cli) SetClusterRouterTable(t string) error {
	return cli.QueryNV("SETCLUSTERROUTERTABLE", toCGPText(t))
}

func (cli *Cli) SetClusterSettings(s map[string]any) error {
	return cli.setMapAny("SETCLUSTERSETTINGS", "", s)
}

func (cli *Cli) SetClusterSignalRules(r []SignalRule) error {
	return setRulesGeneric(cli, "SETCLUSTERSIGNALRULES", "", r)
}

func (cli *Cli) SetClusterWhiteHoleIPs(l IPList) error {
	return cli.setIPListTyped("SETCLUSTERWHITEHOLEIPS", l)
}

func (cli *Cli) SetDebugIPs(l IPList) error  { return cli.setIPListTyped("SETDEBUGIPS", l) }
func (cli *Cli) SetDeniedIPs(l IPList) error { return cli.setIPListTyped("SETDENIEDIPS", l) }

func (cli *Cli) SetDNRSettings(s map[string]any) error {
	return cli.setMapAny("SETDNRSETTINGS", "", s)
}

func (cli *Cli) SetLANIPs(l IPList) error { return cli.setIPListTyped("SETLANIPS", l) }

func (cli *Cli) SetMediaServerSettings(s map[string]any) error {
	return cli.setMapAny("SETMEDIASERVERSETTINGS", "", s)
}

func (cli *Cli) SetModule(n string, s map[string]any) error { return cli.setMapAny("SETMODULE", n, s) }
func (cli *Cli) SetNatedIPs(l IPList) error                 { return cli.setIPListTyped("SETNATEDIPS", l) }
func (cli *Cli) SetNatSiteIPs(l IPList) error               { return cli.setIPListTyped("SETNATSITEIPS", l) }

func (cli *Cli) SetNetwork(s map[string]any) error { return cli.setMapAny("SETNETWORK", "", s) }

func (cli *Cli) SetQueueSettings(s map[string]any) error {
	return cli.setMapAny("SETQUEUESETTINGS", "", s)
}

func (cli *Cli) SetRouterSettings(s map[string]any) error {
	return cli.setMapAny("SETROUTERSETTINGS", "", s)
}

func (cli *Cli) SetRouterTable(table RouterList) error {
	return cli.QueryNV("SETROUTERTABLE", table.MarshalCGP())
}

func (cli *Cli) SetServerIntercept(s map[string]any) error {
	return cli.setMapAny("SETSERVERINTERCEPT", "", s)
}

func (cli *Cli) SetServerMailRules(r []MailRule) error {
	return setRulesGeneric(cli, "SETSERVERMAILRULES", "", r)
}

func (cli *Cli) SetServerSignalRules(r []SignalRule) error {
	return setRulesGeneric(cli, "SETSERVERSIGNALRULES", "", r)
}

func (cli *Cli) SetSessionSettings(s map[string]any) error {
	return cli.setMapAny("SETSESSIONSETTINGS", "", s)
}

func (cli *Cli) SetSignalSettings(s map[string]any) error {
	return cli.setMapAny("SETSIGNALSETTINGS", "", s)
}

func (cli *Cli) SetSMTPSendProfiles(s map[string]any) error {
	return cli.setMapAny("SETSMTPSENDPROFILES", "", s)
}

func (cli *Cli) SetWhiteHoleIPs(l IPList) error { return cli.setIPListTyped("SETWHITEHOLEIPS", l) }

func (cli *Cli) UpdateLogSettings(s map[string]any) error {
	return cli.setMapAny("UPDATELOGSETTINGS", "", s)
}

func (cli *Cli) UpdateModule(n string, s map[string]any) error {
	return cli.setMapAny("UPDATEMODULE", n, s)
}

func (cli *Cli) UpdateServerSettings(s map[string]any) error {
	return cli.setMapAny("UPDATESERVERSETTINGS", "", s)
}

func (cli *Cli) getBannedGeneric(cmd string) (*BannedLines, error) {
	m, err := cli.getMapAny(cmd, "")
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

func (cli *Cli) getIPList(cmd string) (string, error) {
	s, err := cli.getString(cmd)
	if err != nil {
		return "", err
	}
	return toGoText(s), nil
}

func (cli *Cli) getIPListTyped(cmd string) (IPList, error) {
	raw, err := cli.getIPList(cmd)
	if err != nil {
		return nil, err
	}
	return parseIPList(raw), nil
}

func (cli *Cli) setBannedGeneric(cmd string, lines *BannedLines) error {
	if lines == nil {
		return fmt.Errorf("%s: BannedLines is nil", cmd)
	}
	data := map[string]any{
		"BannedBodyLines":   lines.BodyLines,
		"BannedHeaderLines": lines.HeaderLines,
	}
	return cli.setMapAny(cmd, "", data)
}

func (cli *Cli) setIPList(cmd, list string) error {
	return cli.QueryNV(cmd, toCGPText(list))
}

func (cli *Cli) setIPListTyped(cmd string, list IPList) error {
	return cli.setIPList(cmd, list.marshal())
}
