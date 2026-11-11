package netcheck

import "testing"

func mapped(client string, port int) Mapping {
	return Mapping{Known: true, Exists: true, Enabled: true, Client: client, InternalPort: port}
}

func base() Facts {
	return Facts{
		GamePort: 27102, QueryPort: 27131, RconPort: 27015,
		ServerRunning: true, LocalAnswer: true,
		PublicIP: "203.0.113.7", LocalIPs: []string{"192.168.0.10"},
		Router: Router{
			Found: true, ExternalIP: "203.0.113.7",
			Game: mapped("192.168.0.10", 27102), Query: mapped("192.168.0.10", 27131),
			Rcon: Mapping{Known: true},
		},
		Firewall: Yes,
	}
}

func keys(r Result) map[string]string {
	out := map[string]string{}
	for _, it := range r.Items {
		out[it.Key] = it.Status
	}
	return out
}

func TestAllGoodButNoLoopbackIsLikely(t *testing.T) {
	r := Evaluate(base())
	if r.Verdict != VerdictLikely {
		t.Fatalf("verdict %s, items %v", r.Verdict, r.Items)
	}
	k := keys(r)
	if k["routerMapped"] != OK || k["publicOk"] != OK || k["loopbackNo"] != Info || k["localOk"] != OK {
		t.Fatalf("items %v", r.Items)
	}
}

func TestLoopbackAnswerIsOpen(t *testing.T) {
	f := base()
	f.Loopback, f.Steam = true, Yes
	if r := Evaluate(f); r.Verdict != VerdictOpen || keys(r)["steamListed"] != OK {
		t.Fatalf("got %+v", r)
	}
}

func TestCGNATIsAProblem(t *testing.T) {
	f := base()
	f.Router.ExternalIP = "100.72.1.9"
	f.Loopback = true
	r := Evaluate(f)
	if r.Verdict != VerdictProblem || keys(r)["cgnat"] != Bad {
		t.Fatalf("got %+v", r)
	}
}

func TestDifferentPublicIPWarns(t *testing.T) {
	f := base()
	f.PublicIP = "198.51.100.4"
	r := Evaluate(f)
	if keys(r)["ipMismatch"] != Warn || r.Verdict != VerdictUnsure {
		t.Fatalf("got %+v", r)
	}
}

func TestMappingToOtherPCAndWrongPort(t *testing.T) {
	f := base()
	f.Router.Game = mapped("192.168.0.99", 27102)
	f.Router.Query = mapped("192.168.0.10", 27000)
	r := Evaluate(f)
	k := keys(r)
	if k["routerOtherPc"] != Bad || k["routerWrongPort"] != Bad || r.Verdict != VerdictProblem {
		t.Fatalf("got %+v", r)
	}
}

func TestRconForwardedIsAProblem(t *testing.T) {
	f := base()
	f.Rcon = true
	f.Router.Rcon = mapped("192.168.0.10", 27015)
	if r := Evaluate(f); keys(r)["rconForwarded"] != Bad || r.Verdict != VerdictProblem {
		t.Fatalf("got %+v", r)
	}
}

func TestNoRouterNoServer(t *testing.T) {
	r := Evaluate(Facts{GamePort: 1, QueryPort: 2})
	k := keys(r)
	if k["serverOff"] != Info || k["routerUnknown"] != Info || k["publicUnknown"] != Info || r.Verdict != VerdictUnsure {
		t.Fatalf("got %+v", r)
	}
	if _, ok := k["loopbackNo"]; ok {
		t.Fatal("loopback item without a running server")
	}
}

func TestLocalNoAnswerAndMissingMapWarn(t *testing.T) {
	f := base()
	f.LocalAnswer = false
	f.Router.Query = Mapping{Known: true}
	r := Evaluate(f)
	k := keys(r)
	if k["localNo"] != Bad || k["routerNoMap"] != Warn {
		t.Fatalf("got %+v", r)
	}
}

func TestPublicIPMustBeIPv4(t *testing.T) {
	f := base()
	f.PublicIP = "not-an-ip"
	if r := Evaluate(f); r.PublicIP != "" {
		t.Fatalf("public ip %q", r.PublicIP)
	}
}
