package runner

import (
	"bytes"
	"io"
	"os"
)

// Tail lê o que foi acrescentado a um arquivo de log desde a última leitura. Abre e fecha a cada leitura para
// não segurar o arquivo (o servidor renomeia o Insurgency.log para backup ao iniciar).
type Tail struct {
	Path    string
	offset  int64
	partial []byte
}

// maxRead limita o que se lê de uma vez (o log de início tem alguns MB).
const maxRead = 4 << 20

// Reset volta ao começo (arquivo novo).
func (t *Tail) Reset() { t.offset, t.partial = 0, nil }

// SkipToEnd ignora o que já existe (ao acompanhar um servidor que já estava rodando, mostra só o fim).
func (t *Tail) SkipToEnd(keep int64) {
	if info, err := os.Stat(t.Path); err == nil {
		t.offset = info.Size() - keep
		if t.offset < 0 {
			t.offset = 0
		}
		t.partial = nil
	}
}

// Read devolve as linhas completas novas.
func (t *Tail) Read() ([]string, error) {
	f, err := os.Open(t.Path) // #nosec G304 -- log do servidor na pasta configurada
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() < t.offset { // arquivo novo (servidor reiniciado)
		t.Reset()
	}
	if info.Size() == t.offset {
		return nil, nil
	}
	if _, err := f.Seek(t.offset, io.SeekStart); err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(io.LimitReader(f, maxRead))
	if err != nil {
		return nil, err
	}
	t.offset += int64(len(buf))
	data := append(t.partial, buf...)
	last := bytes.LastIndexByte(data, '\n')
	if last < 0 {
		t.partial = data
		return nil, nil
	}
	t.partial = append([]byte(nil), data[last+1:]...)
	var lines []string
	for _, l := range bytes.Split(data[:last], []byte{'\n'}) {
		lines = append(lines, string(bytes.TrimRight(l, "\r")))
	}
	return lines, nil
}
