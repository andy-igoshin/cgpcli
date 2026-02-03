package cgpcli

import (
	"testing"
)

func TestStatistics_Comprehensive(t *testing.T) {
	cli := getTestCli(t)
	defer cli.Close()

	// "*" указывает на текущий аккаунт/домен сессии
	const (
		testAcc    = "*"
		testDomain = "*"
		statKey    = "Logins"
	)

	t.Run("AccountStats", func(t *testing.T) {
		// Get Full
		if _, err := cli.GetAccountStat(testAcc, ""); err != nil {
			t.Errorf("GetAccountStat full failed: %v", err)
		}
		// Get Key
		if _, err := cli.GetAccountStat(testAcc, statKey); err != nil {
			t.Errorf("GetAccountStat key failed: %v", err)
		}
		// Reset Key
		if err := cli.ResetAccountStat(testAcc, statKey); err != nil {
			t.Errorf("ResetAccountStat failed: %v", err)
		}
		if err := cli.ResetAccountStat(testAcc, ""); err != nil {
			t.Errorf("ResetAccountStat failed: %v", err)
		}
	})

	t.Run("DomainStats", func(t *testing.T) {
		// Get Full
		if _, err := cli.GetDomainStat(testDomain, ""); err != nil {
			t.Errorf("GetDomainStat full failed: %v", err)
		}
		if _, err := cli.GetDomainStat(testDomain, statKey); err != nil {
			t.Errorf("GetDomainStat full failed: %v", err)
		}
		if err := cli.ResetDomainStat(testDomain, statKey); err != nil {
			t.Errorf("ResetDomainStat key failed: %v", err)
		}
		if err := cli.ResetDomainStat(testDomain, ""); err != nil {
			t.Errorf("ResetDomainStat full failed: %v", err)
		}
	})
}
