package rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	in := Packet{ID: 7, Type: TypeExec, Body: "scenarios"}
	b, err := Encode(in)
	if err != nil {
		t.Fatal(err)
	}
	if got := binary.LittleEndian.Uint32(b); got != uint32(8+len(in.Body)+2) {
		t.Fatalf("size = %d", got)
	}
	if !bytes.HasSuffix(b, []byte{0, 0}) {
		t.Fatal("missing trailer")
	}
	out, err := Decode(bytes.NewReader(b))
	if err != nil || out != in {
		t.Fatalf("got %+v, %v", out, err)
	}

	neg, _ := Encode(Packet{ID: AuthFailedID, Type: TypeAuthResponse})
	if p, _ := Decode(bytes.NewReader(neg)); p.ID != -1 || p.Body != "" {
		t.Fatalf("auth failed packet = %+v", p)
	}
}

func TestEncodeRefuses(t *testing.T) {
	if _, err := Encode(Packet{Body: "a\x00b"}); !errors.Is(err, ErrBodyNul) {
		t.Fatalf("NUL: %v", err)
	}
	if _, err := Encode(Packet{Body: strings.Repeat("a", MaxBody+1)}); !errors.Is(err, ErrBodyTooLarge) {
		t.Fatalf("large: %v", err)
	}
}

func TestDecodeRefusesBadSize(t *testing.T) {
	for _, size := range []uint32{0, 9, MaxBody + 11, 0xFFFFFFFF} {
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], size)
		if _, err := Decode(bytes.NewReader(b[:])); !errors.Is(err, ErrBadSize) {
			t.Fatalf("size %d: %v", size, err)
		}
	}
	// Pacote cortado no meio.
	b, _ := Encode(Packet{ID: 1, Body: "hello"})
	if _, err := Decode(bytes.NewReader(b[:len(b)-3])); err == nil {
		t.Fatal("truncated packet accepted")
	}
}
