package winsys

import (
	"path/filepath"
	"testing"
)

func TestRuleNames(t *testing.T) {
	progs := ServerPrograms(filepath.Join("C:", "srv"))
	if len(progs) != 2 {
		t.Fatal(progs)
	}
	if got := RuleName(progs[1]); got != RulePrefix+"InsurgencyServer-Win64-Shipping" {
		t.Fatal(got)
	}
	for _, a := range LocalAddresses() {
		if a.IP == "127.0.0.1" {
			t.Fatal("loopback listed")
		}
	}
	if FirewallRule(filepath.Join("C:", "no-such-dir", "nothing.exe")) {
		t.Fatal("rule for a missing program")
	}
}
