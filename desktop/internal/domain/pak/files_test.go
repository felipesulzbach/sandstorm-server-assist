package pak

import (
	"testing"
)

// writeEncoded monta uma entrada codificada sem compressão: campo de bits (posição e tamanho de 32 bits) + posição +
// tamanho.
func writeEncoded(w *buf, size uint32) {
	bits := uint32(1<<31 | 1<<30 | 1<<29)
	w.i32(int32(bits)) // #nosec G115 -- teste
	w.i32(0)
	w.i32(int32(size)) // #nosec G115 -- teste
}

// modernIndex monta o índice principal e o de pastas da versão 11 com dois arquivos codificados e um não codificado.
func modernIndex() (primary, dirIndex []byte) {
	var enc buf
	writeEncoded(&enc, 92328) // posição 0
	writeEncoded(&enc, 1234)  // posição 12
	var p buf
	p.str("../../../")
	p.i32(3)
	p.i64(0)
	p.i32(0) // sem índice por hash
	p.i32(1) // com índice de pastas
	p.i64(0)
	p.i64(1)
	p.b = append(p.b, make([]byte, 20)...)
	p.i32(int32(len(enc.b))) // #nosec G115 -- teste
	p.b = append(p.b, enc.b...)
	p.i32(1) // um arquivo no formato antigo
	p.i64(0)
	p.i64(77)
	p.i64(77)
	p.i32(0)
	p.b = append(p.b, make([]byte, 20+1+4)...)
	var d buf
	d.i32(1)
	d.str("Insurgency/Content/Environment/Props/Modular/Stairs/")
	d.i32(3)
	d.str("SM_Railing.uasset")
	d.i32(12)
	d.str("SM_Railing.uexp")
	d.i32(0)
	d.str("Scenario_Test_Push.uasset")
	d.i32(-1)
	return p.b, d.b
}

func TestModernFiles(t *testing.T) {
	primary, dir := modernIndex()
	files, err := ModernFiles(primary, dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int64{
		"Insurgency/Content/Environment/Props/Modular/Stairs/SM_Railing.uasset":         1234,
		"Insurgency/Content/Environment/Props/Modular/Stairs/SM_Railing.uexp":           92328,
		"Insurgency/Content/Environment/Props/Modular/Stairs/Scenario_Test_Push.uasset": 77,
	}
	if len(files) != len(want) {
		t.Fatalf("files %v", files)
	}
	for _, f := range files {
		if want[f.Path] != f.Size {
			t.Errorf("%s: size %d", f.Path, f.Size)
		}
	}
	if c := ContentsOf(files); len(c.Scenarios) != 1 || c.Scenarios[0] != "Scenario_Test_Push" {
		t.Fatalf("contents %+v", c)
	}
}

func TestModernFilesRejectsBadInput(t *testing.T) {
	primary, dir := modernIndex()
	for i := 0; i < len(primary); i += 7 {
		_, _ = ModernFiles(primary[:i], dir) // nunca pode entrar em pânico
	}
	for i := 0; i < len(dir); i += 5 {
		_, _ = ModernFiles(primary, dir[:i])
	}
	var bad buf
	bad.i32(1)
	bad.str("x/")
	bad.i32(1)
	bad.str("a.uasset")
	bad.i32(9999) // entrada fora do buffer
	if _, err := ModernFiles(primary, bad.b); err == nil {
		t.Fatal("accepted an entry outside the encoded buffer")
	}
}

func legacyIndex(methodSize int) []byte {
	var w buf
	w.str("../../../Insurgency/Plugins/Old/")
	w.i32(2)
	for _, f := range []struct {
		name string
		size int64
	}{{"Content/Mutators/OldMut.uasset", 300}, {"Content/Maps/Old.umap", 5000}} {
		w.str(f.name)
		w.i64(0)
		w.i64(f.size)
		w.i64(f.size)
		if methodSize == 1 {
			w.b = append(w.b, 0)
		} else {
			w.i32(0)
		}
		w.b = append(w.b, make([]byte, 20+1+4)...)
	}
	return w.b
}

func TestLegacyFiles(t *testing.T) {
	for _, ms := range []int{4, 1} {
		files, err := LegacyFiles(legacyIndex(ms))
		if err != nil {
			t.Fatalf("method size %d: %v", ms, err)
		}
		if len(files) != 2 || files[0].Path != "Insurgency/Plugins/Old/Content/Mutators/OldMut.uasset" ||
			files[1].Size != 5000 {
			t.Fatalf("method size %d: %+v", ms, files)
		}
		c := ContentsOf(files)
		if len(c.Mutators) != 1 || c.Mutators[0] != "OldMut" || len(c.Maps) != 1 || c.Maps[0] != "Old" {
			t.Fatalf("contents %+v", c)
		}
	}
	idx := legacyIndex(4)
	for i := 0; i < len(idx); i += 3 {
		_, _ = LegacyFiles(idx[:i])
	}
}

func TestFooterAnyAcceptsVersion9(t *testing.T) {
	var f buf
	f.b = append(f.b, make([]byte, 100)...)
	f.b = append(f.b, make([]byte, 17)...)
	f.i32(Magic)
	f.i32(9)
	f.i64(10)
	f.i64(20)
	f.b = append(f.b, make([]byte, 20)...)
	if _, err := ParseFooter(f.b, int64(len(f.b))); err == nil {
		t.Fatal("ParseFooter should keep refusing version 9")
	}
	ft, err := FooterAny(f.b, int64(len(f.b)))
	if err != nil || ft.Version != 9 || ft.IndexOffset != 10 || ft.IndexSize != 20 {
		t.Fatalf("footer %+v %v", ft, err)
	}
}

func FuzzModernFiles(f *testing.F) {
	p, d := modernIndex()
	f.Add(p, d)
	f.Fuzz(func(t *testing.T, primary, dir []byte) {
		_, _ = ModernFiles(primary, dir)
	})
}

func FuzzLegacyFiles(f *testing.F) {
	f.Add(legacyIndex(4))
	f.Fuzz(func(t *testing.T, idx []byte) {
		_, _ = LegacyFiles(idx)
	})
}
