package steamcmd

import (
	"slices"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		line string
		want Line
	}{
		{" Update state (0x61) downloading, progress: 45.23 (1234567890 / 2730000000)", Line{Phase: Downloading, Percent: 45.23, Done: 1234567890, Total: 2730000000}},
		{" Update state (0x5) verifying install, progress: 3.10 (100 / 3226)", Line{Phase: Verifying, Percent: 3.10, Done: 100, Total: 3226}},
		{" Update state (0x11) preallocating, progress: 70.00 (7 / 10)", Line{Phase: Preallocate, Percent: 70, Done: 7, Total: 10}},
		{" Update state (0x101) committing, progress: 99.99 (9 / 10)", Line{Phase: Committing, Percent: 99.99, Done: 9, Total: 10}},
		{"Success! App '581330' fully installed.", Line{Phase: Done, Percent: 100}},
		{"Success! App '581330' already up to date.", Line{Phase: Done, Percent: 100, UpToDate: true}},
		{"ERROR! Failed to install app '581330' (Missing configuration)", Line{Phase: Failed, Error: "Missing configuration"}},
		{"ERROR! Failed to install app '581330' (Disk write failure)", Line{Phase: Failed, Error: "Disk write failure"}},
		{"Logging in user 'anonymous' to Steam Public...OK", Line{Phase: LoggingIn}},
		{"[ 40%] Downloading update (12,345 of 30,000 KB)...", Line{Phase: SelfUpdate, Percent: 40}},
	}
	for _, c := range cases {
		got, ok := Parse(c.line)
		if !ok || got != c.want {
			t.Errorf("Parse(%q) = %+v, %v; want %+v", c.line, got, ok, c.want)
		}
	}
}

func TestArgs(t *testing.T) {
	got := Args(`D:\srv`, 581330, true)
	want := []string{"+force_install_dir", `D:\srv`, "+login", "anonymous", "+app_update", "581330", "validate", "+quit"}
	if !slices.Equal(got, want) {
		t.Fatalf("%v", got)
	}
}

func TestNestedDirs(t *testing.T) {
	cases := []struct {
		server, steamcmd string
		want             bool
	}{
		{`D:\Servers`, `D:\Servers\steamcmd`, true},
		{`D:\Servers\`, `d:\servers\SteamCMD\`, true},
		{`D:\Servers\steamcmd\x`, `D:\Servers\steamcmd`, true},
		{`D:\Servers\Sandstorm`, `D:\Servers\Sandstorm`, true},
		{`D:\Servers\Sandstorm`, `D:\Servers\steamcmd`, false},
		{`D:\Servers\Sand`, `D:\Servers\Sandstorm\steamcmd`, false},
		{`D:\Servers`, ``, false},
	}
	for _, c := range cases {
		if got := NestedDirs(c.server, c.steamcmd); got != c.want {
			t.Errorf("NestedDirs(%q, %q) = %v", c.server, c.steamcmd, got)
		}
	}
}

func TestWrongInstallDir(t *testing.T) {
	if !WrongInstallDir("[2026-11-13 18:00:18] Please set the game install path to something other than the Steam install folder") {
		t.Fatal("did not see the warning")
	}
	if WrongInstallDir("Success! App '581330' fully installed.") {
		t.Fatal("false positive")
	}
}
