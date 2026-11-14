package desktop

import (
	"net"
	"testing"
)

// Um programa só em 127.0.0.1 (como o MongoDB na 27017) também conta como porta ocupada.
func TestPortBusyLoopbackOnly(t *testing.T) {
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	if !portBusy(port, true) {
		t.Error("TCP port used on 127.0.0.1 seen as free")
	}
	_ = l.Close()
	if portBusy(port, true) {
		t.Error("closed TCP port seen as busy")
	}

	u, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	uport := u.LocalAddr().(*net.UDPAddr).Port
	if !portBusy(uport, false) {
		t.Error("UDP port used on 127.0.0.1 seen as free")
	}
	_ = u.Close()
}
