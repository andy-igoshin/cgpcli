package cgpcli

import (
	"fmt"
	"testing"
	"time"
)

// --- СЕКЦИЯ 1: IP LISTS (Локальные и Кластерные) ---

func TestServer_IPLists_Honest(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Список пар (Название, Getter, Setter)
	localIPs := []struct {
		name string
		get  func() (IPList, error)
		set  func(IPList) error
	}{
		{"Blacklisted", cli.GetBlacklistedIPs, cli.SetBlacklistedIPs},
		{"Client", cli.GetClientIPs, cli.SetClientIPs},
		{"Debug", cli.GetDebugIPs, cli.SetDebugIPs},
		{"Denied", cli.GetDeniedIPs, cli.SetDeniedIPs},
		{"LAN", cli.GetLANIPs, cli.SetLANIPs},
		{"Nated", cli.GetNatedIPs, cli.SetNatedIPs},
		{"NatSite", cli.GetNatSiteIPs, cli.SetNatSiteIPs},
		{"WhiteHole", cli.GetWhiteHoleIPs, cli.SetWhiteHoleIPs},
	}

	for _, tt := range localIPs {
		t.Run("Local_"+tt.name, func(t *testing.T) {
			orig, err := tt.get()
			if err != nil {
				t.Fatalf("GET %s failed: %v", tt.name, err)
			}
			// Round-trip: проверяем, что маршалер не портит структуру при записи
			if err := tt.set(orig); err != nil {
				t.Errorf("SET %s failed: %v", tt.name, err)
			}
		})
	}

	// Кластерные IP-списки
	clusterIPs := []struct {
		name string
		get  func() (IPList, error)
		set  func(IPList) error
	}{
		{"ClusterBlacklisted", cli.GetClusterBlacklistedIPs, cli.SetClusterBlacklistedIPs},
		{"ClusterClient", cli.GetClusterClientIPs, cli.SetClusterClientIPs},
		{"ClusterDebug", cli.GetClusterDebugIPs, cli.SetClusterDebugIPs},
		{"ClusterDenied", cli.GetClusterDeniedIPs, cli.SetClusterDeniedIPs},
		{"ClusterLAN", cli.GetClusterLANIPs, cli.SetClusterLANIPs},
		{"ClusterNated", cli.GetClusterNatedIPs, cli.SetClusterNatedIPs},
		{"ClusterNatSite", cli.GetClusterNatSiteIPs, cli.SetClusterNatSiteIPs},
		{"ClusterWhiteHole", cli.GetClusterWhiteHoleIPs, cli.SetClusterWhiteHoleIPs},
	}

	for _, tt := range clusterIPs {
		t.Run(tt.name, func(t *testing.T) {
			orig, err := tt.get()
			if err != nil {
				skipIfNoClusterAuth(t, err)
				t.Fatalf("GET %s failed: %v", tt.name, err)
			}
			if err := tt.set(orig); err != nil {
				t.Errorf("SET %s failed: %v", tt.name, err)
			}
		})
	}
}

// --- СЕКЦИЯ 2: ROUTER TABLES ---

func TestServer_RouterTables_Honest(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	t.Run("LocalRouterTable", func(t *testing.T) {
		orig, err := cli.GetRouterTable()
		if err != nil {
			t.Fatalf("GetRouterTable failed: %v", err)
		}

		testRule := fmt.Sprintf("sdk-test-%d.local = trash", time.Now().Unix())
		workTable := make(RouterList, len(orig))
		copy(workTable, orig)
		workTable.Add(testRule, "SDK RoundTrip Test", 40)

		if err := cli.SetRouterTable(workTable); err != nil {
			t.Fatalf("SetRouterTable failed: %v", err)
		}

		// Сверяем результат
		updated, _ := cli.GetRouterTable()
		_ = cli.SetRouterTable(orig) // Сразу восстанавливаем

		found := false
		for _, r := range updated {
			if r.Rule == testRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Rule %q not found in updated table", testRule)
		}
	})

	t.Run("ClusterRouterTable", func(t *testing.T) {
		orig, err := cli.GetClusterRouterTable()
		if err != nil {
			skipIfNoClusterAuth(t, err)
			t.Fatalf("GetClusterRouterTable failed: %v", err)
		}
		if err := cli.SetClusterRouterTable(orig); err != nil {
			t.Errorf("SetClusterRouterTable failed: %v", err)
		}
	})
}

// --- СЕКЦИЯ 3: SYSTEM SETTINGS (Maps) ---

func TestServer_SystemSettings_RoundTrip(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Локальные мапы
	localMaps := []struct {
		name string
		get  func() (map[string]any, error)
		set  func(map[string]any) error
	}{
		{"DNR", cli.GetDNRSettings, cli.SetDNRSettings},
		{"Log", cli.GetLogSettings, cli.UpdateLogSettings},
		{"Network", cli.GetNetwork, cli.SetNetwork},
		{"Queue", cli.GetQueueSettings, cli.SetQueueSettings},
		{"Router", cli.GetRouterSettings, cli.SetRouterSettings},
		{"ServerSettings", cli.GetServerSettings, cli.UpdateServerSettings},
		{"Session", cli.GetSessionSettings, cli.SetSessionSettings},
		{"Signal", cli.GetSignalSettings, cli.SetSignalSettings},
		{"MediaServer", cli.GetMediaServerSettings, cli.SetMediaServerSettings},
		{"ServerIntercept", cli.GetServerIntercept, cli.SetServerIntercept},
		{"SMTPSendProfiles", cli.GetSMTPSendProfiles, cli.SetSMTPSendProfiles},
	}

	for _, tt := range localMaps {
		t.Run("Local_"+tt.name, func(t *testing.T) {
			orig, err := tt.get()
			if err != nil {
				t.Fatalf("GET %s failed: %v", tt.name, err)
			}
			if err := tt.set(orig); err != nil {
				t.Errorf("SET/UPDATE %s failed: %v", tt.name, err)
			}
		})
	}

	// Кластерные мапы
	clusterMaps := []struct {
		name string
		get  func() (map[string]any, error)
		set  func(map[string]any) error
	}{
		{"ClusterIntercept", cli.GetClusterIntercept, cli.SetClusterIntercept},
		{"ClusterNetwork", cli.GetClusterNetwork, cli.SetClusterNetwork},
		{"ClusterRouter", cli.GetClusterRouterSettings, cli.SetClusterRouterSettings},
		{"ClusterSettings", cli.GetClusterSettings, cli.SetClusterSettings},
	}

	for _, tt := range clusterMaps {
		t.Run(tt.name, func(t *testing.T) {
			orig, err := tt.get()
			if err != nil {
				skipIfNoClusterAuth(t, err)
				t.Fatalf("GET %s failed: %v", tt.name, err)
			}
			if err := tt.set(orig); err != nil {
				skipIfNoClusterAuth(t, err)
				t.Errorf("SET %s failed: %v", tt.name, err)
			}
		})
	}
}

// --- СЕКЦИЯ 4: RULES & BANNED ---

func TestServer_RulesAndBanned_Honest(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// Banned Lines
	t.Run("BannedLines_Local", func(t *testing.T) {
		orig, err := cli.GetBanned()
		if err != nil {
			t.Fatalf("GetBanned failed: %v", err)
		}
		if err := cli.SetBanned(orig); err != nil {
			t.Errorf("SetBanned failed: %v", err)
		}
	})

	t.Run("BannedLines_Cluster", func(t *testing.T) {
		orig, err := cli.GetClusterBanned()
		if err != nil {
			skipIfNoClusterAuth(t, err)
		}
		if err == nil {
			_ = cli.SetClusterBanned(orig)
		}
	})

	// Mail Rules
	t.Run("MailRules_Local", func(t *testing.T) {
		orig, err := cli.GetServerMailRules()
		if err != nil {
			t.Fatalf("GetServerMailRules failed: %v", err)
		}
		if err := cli.SetServerMailRules(orig); err != nil {
			t.Errorf("SetServerMailRules failed: %v", err)
		}
	})

	t.Run("MailRules_Cluster", func(t *testing.T) {
		orig, err := cli.GetClusterMailRules()
		if err != nil {
			skipIfNoClusterAuth(t, err)
		}
		if err == nil {
			_ = cli.SetClusterMailRules(orig)
		}
	})

	// Signal Rules
	t.Run("SignalRules_Local", func(t *testing.T) {
		orig, err := cli.GetServerSignalRules()
		if err != nil {
			t.Fatalf("GetServerSignalRules failed: %v", err)
		}
		if err := cli.SetServerSignalRules(orig); err != nil {
			t.Errorf("SetServerSignalRules failed: %v", err)
		}
	})

	t.Run("SignalRules_Cluster", func(t *testing.T) {
		orig, err := cli.GetClusterSignalRules()
		if err != nil {
			skipIfNoClusterAuth(t, err)
		}
		if err == nil {
			_ = cli.SetClusterSignalRules(orig)
		}
	})
}

// --- СЕКЦИЯ 5: UTILITIES & MODULES ---

func TestServer_UtilityMethods(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	t.Run("GetIPState", func(t *testing.T) {
		if _, err := cli.GetIPState("127.0.0.1", false); err != nil {
			t.Errorf("GetIPState failed: %v", err)
		}
	})

	t.Run("Route_recursive", func(t *testing.T) {
		if _, err := cli.Route("postmaster", "mail"); err != nil {
			t.Errorf("Route failed: %v", err)
		}
	})

	t.Run("ModuleOperations", func(t *testing.T) {
		mods, err := cli.ListModules()
		if err != nil {
			t.Fatalf("ListModules failed: %v", err)
		}
		if len(mods) > 0 {
			mName := mods[0]
			cfg, err := cli.GetModule(mName)
			if err != nil {
				t.Fatalf("GetModule(%s) failed: %v", mName, err)
			}

			// Проверка Set и Update
			if err := cli.SetModule(mName, cfg); err != nil {
				t.Errorf("SetModule failed: %v", err)
			}
			if err := cli.UpdateModule(mName, cfg); err != nil {
				t.Errorf("UpdateModule failed: %v", err)
			}
		}
	})

	t.Run("RefreshOS", func(t *testing.T) {
		if err := cli.RefreshOSData(); err != nil {
			t.Errorf("RefreshOSData failed: %v", err)
		}
	})
}
