package pak

import (
	"encoding/binary"
	"testing"
)

type buf struct{ b []byte }

func (w *buf) i32(v int32) { w.b = binary.LittleEndian.AppendUint32(w.b, uint32(v)) } // #nosec G115 -- teste
func (w *buf) i64(v int64) { w.b = binary.LittleEndian.AppendUint64(w.b, uint64(v)) } // #nosec G115 -- teste
func (w *buf) str(s string) {
	w.i32(int32(len(s) + 1)) // #nosec G115 -- teste
	w.b = append(append(w.b, s...), 0)
}

// fakePak monta um .pak mínimo: dados, índice de pastas, índice principal e rodapé (versão 11).
func fakePak(dirs map[string][]string) []byte {
	var f buf
	f.b = append(f.b, make([]byte, 64)...) // "dados"
	dirOff := int64(len(f.b))
	var d buf
	d.i32(int32(len(dirs))) // #nosec G115 -- teste
	for dir, files := range dirs {
		d.str(dir)
		d.i32(int32(len(files))) // #nosec G115 -- teste
		for i, name := range files {
			d.str(name)
			d.i32(int32(i)) // #nosec G115 -- teste
		}
	}
	f.b = append(f.b, d.b...)
	idxOff := int64(len(f.b))
	var p buf
	p.str("../../../Insurgency/Plugins/Test/")
	p.i32(3)
	p.i64(0)
	p.i32(1) // tem índice por hash
	p.i64(0)
	p.i64(0)
	p.b = append(p.b, make([]byte, 20)...)
	p.i32(1) // tem índice de pastas
	p.i64(dirOff)
	p.i64(int64(len(d.b)))
	f.b = append(f.b, p.b...)
	f.b = append(f.b, make([]byte, 17)...) // GUID da chave + "não cifrado"
	f.i32(Magic)
	f.i32(11)
	f.i64(idxOff)
	f.i64(int64(len(p.b)))
	f.b = append(f.b, make([]byte, 160)...) // hash e nomes de compressão
	return f.b
}

func scan(t *testing.T, file []byte) []string {
	t.Helper()
	size := int64(len(file))
	tail := file[max(0, len(file)-TailSize):]
	ft, err := ParseFooter(tail, size)
	if err != nil {
		t.Fatal(err)
	}
	off, n, err := DirIndex(file[ft.IndexOffset:ft.IndexOffset+ft.IndexSize], size)
	if err != nil {
		t.Fatal(err)
	}
	names, err := Mutators(file[off : off+n])
	if err != nil {
		t.Fatal(err)
	}
	return names
}

func TestMutators(t *testing.T) {
	file := fakePak(map[string][]string{
		"Content/Mutators/":  {"NoRestrictedArea.uasset", "NoRestrictedArea.uexp", "Bad Name.uasset"},
		"Content/Maps/":      {"Mutators.umap"},
		"Content/Mutators2/": {"Other.uasset"},
	})
	got := scan(t, file)
	if len(got) != 1 || got[0] != "NoRestrictedArea" {
		t.Fatalf("got %v", got)
	}
	if got := scan(t, fakePak(map[string][]string{"Content/Maps/": {"Town.umap"}})); len(got) != 0 {
		t.Fatalf("map mod: %v", got)
	}
}

func TestBadInput(t *testing.T) {
	file := fakePak(map[string][]string{"Content/Mutators/": {"X.uasset"}})
	// Rodapé apontando para fora do arquivo.
	if _, err := ParseFooter(file[max(0, len(file)-TailSize):], 100); err == nil {
		t.Fatal("index outside the file accepted")
	}
	// Índice cifrado.
	enc := append([]byte(nil), file...)
	i := len(enc) - 160 - 24 - 1
	enc[i] = 1
	if _, err := ParseFooter(enc[max(0, len(enc)-TailSize):], int64(len(enc))); err == nil {
		t.Fatal("encrypted index accepted")
	}
	// Dados cortados e contagens absurdas não podem estourar.
	if _, err := Mutators([]byte{1, 0, 0}); err == nil {
		t.Fatal("truncated index accepted")
	}
	if _, err := Mutators([]byte{0xff, 0xff, 0xff, 0x7f}); err == nil {
		t.Fatal("huge count accepted")
	}
	if _, _, err := DirIndex([]byte{0xff, 0xff, 0xff, 0x7f}, 10); err == nil {
		t.Fatal("huge string accepted")
	}
}

func TestPluginName(t *testing.T) {
	blob := "PluginName:NoRestrictedArea,ContentDirectory:../../../Insurgency/Plugins/NoRestrictedArea/Content"
	if got := PluginName(blob); got != "NoRestrictedArea" {
		t.Fatal(got)
	}
	if got := PluginName("PluginName:bad name;rm,x"); got != "" {
		t.Fatal(got)
	}
}

func FuzzMutators(f *testing.F) {
	f.Add(fakePak(map[string][]string{"Content/Mutators/": {"A.uasset"}}))
	f.Fuzz(func(_ *testing.T, b []byte) {
		_, _ = Mutators(b)
		_, _, _ = DirIndex(b, int64(len(b)))
		_, _ = ParseFooter(b, int64(len(b)))
	})
}

func TestScanMapsAndScenarios(t *testing.T) {
	file := fakePak(map[string][]string{
		"Content/Maps/":      {"Old_World.umap", "Old_World_AI.umap", "Old_World_Checkpoint_Security.umap"},
		"Content/Scenarios/": {"Scenario_Old_World_Checkpoint_Security_Hardcore.uasset", "Scenario_Old_World_Push.uexp"},
		"Content/Mutators/":  {"Quiet.uasset"},
	})
	size := int64(len(file))
	ft, _ := ParseFooter(file[max(0, len(file)-TailSize):], size)
	off, n, _ := DirIndex(file[ft.IndexOffset:ft.IndexOffset+ft.IndexSize], size)
	c, err := Scan(file[off : off+n])
	if err != nil || len(c.Maps) != 3 || len(c.Scenarios) != 1 || len(c.Mutators) != 1 {
		t.Fatalf("%+v %v", c, err)
	}
	if got := ScenarioMap(c.Scenarios[0], c.Maps); got != "Old_World" {
		t.Fatal(got)
	}
	if got := ScenarioMap("Scenario_Town_Checkpoint", c.Maps); got != "" {
		t.Fatal(got)
	}
}
