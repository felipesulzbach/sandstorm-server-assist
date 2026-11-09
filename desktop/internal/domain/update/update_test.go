package update

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

const testSHA = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func signed(t *testing.T, m Manifest) (Manifest, ed25519.PublicKey) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	m.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, m.Message()))
	return m, pub
}

func TestVerify(t *testing.T) {
	good := Manifest{Version: "0.2.0", File: InstallerName("0.2.0"), SHA256: testSHA, Size: 1000}
	m, pub := signed(t, good)
	if err := m.Verify(pub); err != nil {
		t.Fatalf("valid manifest refused: %v", err)
	}
	if m.File != "sandstorm-server-assist-0.2.0-windows-setup.exe" {
		t.Fatalf("installer name: %s", m.File)
	}

	// Mudar qualquer campo assinado quebra a assinatura.
	tampered := m
	tampered.Size = 999
	if tampered.Verify(pub) != ErrBadSignature {
		t.Fatal("tampered size accepted")
	}
	other, _ := signed(t, good)
	if other.Verify(pub) != ErrBadSignature {
		t.Fatal("signature from another key accepted")
	}

	bad := []Manifest{
		{Version: "0.2", File: InstallerName("0.2"), SHA256: testSHA, Size: 1},
		{Version: "0.2.0", File: InstallerName("0.3.0"), SHA256: testSHA, Size: 1},
		{Version: "0.2.0", File: "../evil.exe", SHA256: testSHA, Size: 1},
		{Version: "0.2.0", File: InstallerName("0.2.0"), SHA256: "XYZ", Size: 1},
		{Version: "0.2.0", File: InstallerName("0.2.0"), SHA256: testSHA, Size: 0},
		{Version: "0.2.0", File: InstallerName("0.2.0"), SHA256: testSHA, Size: MaxInstallerBytes + 1},
	}
	for i, b := range bad {
		s, k := signed(t, b)
		if s.Verify(k) != ErrBadManifest {
			t.Errorf("case %d: invalid manifest accepted", i)
		}
	}
}

func TestNewer(t *testing.T) {
	cases := []struct {
		cur, cand string
		want      bool
	}{
		{"0.1.0", "0.1.1", true},
		{"0.1.0", "0.2.0", true},
		{"0.9.9", "1.0.0", true},
		{"0.1.10", "0.1.9", false},
		{"0.2.0", "0.2.0", false},
		{"0.1.0-dev", "9.9.9", false},
		{"0.1.0", "x", false},
	}
	for _, c := range cases {
		if got := Newer(c.cur, c.cand); got != c.want {
			t.Errorf("Newer(%q,%q)=%v", c.cur, c.cand, got)
		}
	}
	if IsRelease("0.1.0-dev") || !IsRelease("0.1.0") {
		t.Fatal("IsRelease")
	}
}
