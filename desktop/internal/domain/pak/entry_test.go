package pak

import (
	"bytes"
	"compress/zlib"
	"testing"
)

// entryRecord monta cabeçalho + dados de um arquivo (comprimido em Zlib num bloco, ou sem compressão).
func entryRecord(t *testing.T, data []byte, compress bool) []byte {
	t.Helper()
	payload := data
	if compress {
		var z bytes.Buffer
		w := zlib.NewWriter(&z)
		_, _ = w.Write(data)
		_ = w.Close()
		payload = z.Bytes()
	}
	var h buf
	h.i64(0)
	h.i64(int64(len(payload)))
	h.i64(int64(len(data)))
	if compress {
		h.i32(1)
	} else {
		h.i32(0)
	}
	h.b = append(h.b, make([]byte, 20)...)
	headerSize := int64(len(h.b) + 1 + 4)
	if compress {
		headerSize += 4 + 16
		h.i32(1)
		h.i64(headerSize)
		h.i64(headerSize + int64(len(payload)))
	}
	h.b = append(h.b, 0)
	h.i32(65536)
	return append(h.b, payload...)
}

func readFrom(rec []byte) func(pos, n int64) ([]byte, error) {
	return func(pos, n int64) ([]byte, error) {
		if pos < 0 || n < 0 || pos+n > int64(len(rec)) {
			return nil, ErrFormat
		}
		return rec[pos : pos+n], nil
	}
}

const fakeIni = "[/Script/Other]\nApiKey=nope\n\n[/Script/Modio.ModioSettings]\nGameId=254\n" +
	"ApiKey=0123456789abcdef0123456789abcdef\nApiRoot=https://example.invalid\n"

func TestReadDataZlibAndPlain(t *testing.T) {
	for _, compress := range []bool{true, false} {
		rec := entryRecord(t, []byte(fakeIni), compress)
		e, err := ParseEntry(rec)
		if err != nil {
			t.Fatal(err)
		}
		data, err := ReadData(e, []string{"Zlib"}, readFrom(rec))
		if err != nil || string(data) != fakeIni {
			t.Fatalf("compress=%v: %q %v", compress, data, err)
		}
		if key, ok := ModioAPIKey(string(data)); !ok || key != "0123456789abcdef0123456789abcdef" {
			t.Fatalf("key %q %v", key, ok)
		}
	}
}

func TestReadDataRefuses(t *testing.T) {
	rec := entryRecord(t, []byte(fakeIni), true)
	e, _ := ParseEntry(rec)
	if _, err := ReadData(e, []string{"Oodle"}, readFrom(rec)); err != ErrUnsupported {
		t.Fatalf("other compression: %v", err)
	}
	enc := e
	enc.Encrypted = true
	if _, err := ReadData(enc, []string{"Zlib"}, readFrom(rec)); err != ErrUnsupported {
		t.Fatalf("encrypted: %v", err)
	}
	bad := e
	bad.UncompressedSize++
	if _, err := ReadData(bad, []string{"Zlib"}, readFrom(rec)); err == nil {
		t.Fatal("size mismatch accepted")
	}
	for i := 0; i < len(rec); i += 3 {
		_, _ = ParseEntry(rec[:i]) // nunca entra em pânico
	}
}

func TestModioAPIKeyNeedsTheGame(t *testing.T) {
	for _, s := range []string{
		"[/Script/Modio.ModioSettings]\nGameId=255\nApiKey=0123456789abcdef0123456789abcdef\n",
		"[/Script/Modio.ModioSettings]\nGameId=254\nApiKey=short\n",
		"[/Script/Other]\nGameId=254\nApiKey=0123456789abcdef0123456789abcdef\n",
	} {
		if _, ok := ModioAPIKey(s); ok {
			t.Errorf("accepted %q", s)
		}
	}
}

func TestMethods(t *testing.T) {
	var f buf
	f.b = append(f.b, make([]byte, 17)...)
	f.i32(Magic)
	f.i32(11)
	f.i64(0)
	f.i64(1)
	f.b = append(f.b, make([]byte, 20)...)
	name := make([]byte, 32)
	copy(name, "Zlib")
	f.b = append(f.b, name...)
	f.b = append(f.b, make([]byte, 32*4)...)
	if m := Methods(f.b); len(m) != 5 || m[0] != "Zlib" || m[1] != "" {
		t.Fatalf("methods %q", m)
	}
}

func FuzzParseEntry(f *testing.F) {
	var seed buf
	seed.i64(0)
	seed.i64(10)
	seed.i64(10)
	seed.i32(0)
	seed.b = append(seed.b, make([]byte, 25)...)
	f.Add(seed.b)
	f.Fuzz(func(t *testing.T, b []byte) {
		e, err := ParseEntry(b)
		if err == nil {
			_, _ = ReadData(e, []string{"Zlib"}, readFrom(b))
		}
	})
}
