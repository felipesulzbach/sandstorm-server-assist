// Package modcheck são as regras para achar mods com problema (ADR-039), a partir dos arquivos que cada mod traz
// (lidos dos índices dos .pak):
//   - Conflicts: dois mods com o mesmo arquivo do jogo em tamanhos diferentes trazem versões diferentes dele; o
//     jogo pode fechar com "Serial size mismatch" ao carregar um mapa.
//   - Assets e Owners: a mensagem de um crash ("Serial size mismatch" de um jogador ou do servidor) cita o objeto;
//     o app acha quais mods instalados trazem esse arquivo.
//
// Nada aqui faz I/O.
package modcheck

import (
	"path"
	"regexp"
	"sort"
	"strings"
)

// File é um arquivo de um mod (caminho a partir da pasta do jogo e tamanho sem compressão).
type File struct {
	Path string
	Size int64
}

// Mod é um mod instalado e os arquivos dos .pak dele.
type Mod struct {
	ID    int
	Name  string
	Files []File
}

// ModRef identifica um mod na tela.
type ModRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Conflict é um par de mods com arquivos iguais em versões diferentes.
type Conflict struct {
	A ModRef `json:"a"`
	B ModRef `json:"b"`
	// Count é quantos objetos (asset) diferem; Example é um deles, no formato do jogo (/Game/...).
	Count   int    `json:"count"`
	Example string `json:"example"`
}

// Extensões que formam um objeto do jogo (o .uasset é o cabeçalho; o .uexp e o .ubulk, os dados).
var assetExt = map[string]bool{".uasset": true, ".uexp": true, ".ubulk": true, ".umap": true}

// stem tira a extensão e deixa em minúsculas: os arquivos de um mesmo objeto ficam juntos.
func stem(p string) (string, bool) {
	ext := strings.ToLower(path.Ext(p))
	if !assetExt[ext] {
		return "", false
	}
	return strings.ToLower(strings.TrimSuffix(p, path.Ext(p))), true
}

// Conflicts compara os mods dois a dois: o mesmo objeto (caminho sem extensão) com algum arquivo de tamanho
// diferente conta como conflito. Cópias idênticas (mesmos tamanhos), comuns entre mods de mapa, não contam.
// Devolve os pares do maior para o menor número de objetos.
func Conflicts(mods []Mod) []Conflict {
	type owner struct {
		mod   int // índice em mods
		sizes map[string]int64
	}
	byAsset := map[string][]owner{}
	display := map[string]string{} // caminho original (com maiúsculas) para mostrar
	for i, m := range mods {
		local := map[string]map[string]int64{}
		for _, f := range m.Files {
			s, ok := stem(f.Path)
			if !ok {
				continue
			}
			if _, seen := display[s]; !seen {
				display[s] = strings.TrimSuffix(f.Path, path.Ext(f.Path))
			}
			if local[s] == nil {
				local[s] = map[string]int64{}
			}
			local[s][strings.ToLower(path.Ext(f.Path))] = f.Size
		}
		for s, sizes := range local {
			byAsset[s] = append(byAsset[s], owner{mod: i, sizes: sizes})
		}
	}
	type key struct{ a, b int }
	found := map[key]*Conflict{}
	for s, owners := range byAsset {
		if len(owners) < 2 {
			continue
		}
		for i := 0; i < len(owners); i++ {
			for j := i + 1; j < len(owners); j++ {
				a, b := owners[i], owners[j]
				if a.mod == b.mod || sameSizes(a.sizes, b.sizes) {
					continue
				}
				x, y := a.mod, b.mod
				if mods[x].ID > mods[y].ID {
					x, y = y, x
				}
				k := key{x, y}
				c := found[k]
				if c == nil {
					c = &Conflict{A: ref(mods[x]), B: ref(mods[y])}
					found[k] = c
				}
				c.Count++
				if ex := GamePath(display[s]); c.Example == "" || ex < c.Example {
					c.Example = ex
				}
			}
		}
	}
	out := make([]Conflict, 0, len(found))
	for _, c := range found {
		out = append(out, *c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].A.ID*1_000_000_000+out[i].B.ID < out[j].A.ID*1_000_000_000+out[j].B.ID
	})
	return out
}

func sameSizes(a, b map[string]int64) bool {
	for ext, sa := range a {
		if sb, ok := b[ext]; ok && sb != sa {
			return false
		}
	}
	return true
}

func ref(m Mod) ModRef { return ModRef{ID: m.ID, Name: m.Name} }

// GamePath converte o caminho de um arquivo (Insurgency/Content/X) no nome do objeto no jogo (/Game/X), para
// mostrar. Outros caminhos voltam como estão.
func GamePath(stemPath string) string {
	low := strings.ToLower(stemPath)
	switch {
	case strings.HasPrefix(low, "insurgency/content/"):
		return "/Game/" + stemPath[len("insurgency/content/"):]
	case strings.HasPrefix(low, "engine/content/"):
		return "/Engine/" + stemPath[len("engine/content/"):]
	}
	return "/" + stemPath
}

// reObject acha nomes de objeto do Unreal numa mensagem: /Raiz/Pasta/Nome ou /Raiz/Pasta/Nome.Nome.
var reObject = regexp.MustCompile(`(?:^|[\s'"(:=])(/[A-Za-z0-9_]+(?:/[A-Za-z0-9_\-]+)+)(?:\.[A-Za-z0-9_\-]+)?`)

// maxObject é o maior nome de objeto aceito (caminhos do jogo são bem menores).
const maxObject = 512

// MaxText limita a mensagem colada.
const MaxText = 64 << 10

// Assets lê os objetos citados numa mensagem de erro (até 20, sem repetir). Ignora /Script (código, não arquivo).
func Assets(text string) []string {
	if len(text) > MaxText {
		text = text[:MaxText]
	}
	seen := map[string]bool{}
	var out []string
	for _, m := range reObject.FindAllStringSubmatch(text, -1) {
		obj := m[1]
		if len(obj) > maxObject {
			continue
		}
		root := strings.ToLower(strings.SplitN(obj[1:], "/", 2)[0])
		if root == "script" || seen[strings.ToLower(obj)] {
			continue
		}
		seen[strings.ToLower(obj)] = true
		out = append(out, obj)
		if len(out) == 20 {
			break
		}
	}
	return out
}

// candidates são os caminhos (sem extensão, minúsculas) onde o objeto pode estar: /Game → Insurgency/Content,
// /Engine → Engine/Content, /<Plugin>/X → qualquer pasta <Plugin>/Content/X (plugins e mods).
func candidates(object string) (exact string, suffix string) {
	parts := strings.SplitN(strings.TrimPrefix(object, "/"), "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	root, rest := strings.ToLower(parts[0]), strings.ToLower(parts[1])
	switch root {
	case "game":
		return "insurgency/content/" + rest, ""
	case "engine":
		return "engine/content/" + rest, ""
	}
	return "", "/" + root + "/content/" + rest
}

// Owners devolve os mods que trazem o objeto (em qualquer das extensões de objeto).
func Owners(object string, mods []Mod) []ModRef {
	exact, suffix := candidates(object)
	if exact == "" && suffix == "" {
		return nil
	}
	var out []ModRef
	for _, m := range mods {
		for _, f := range m.Files {
			s, ok := stem(f.Path)
			if !ok {
				continue
			}
			if (exact != "" && s == exact) || (suffix != "" && strings.HasSuffix("/"+s, suffix)) {
				out = append(out, ref(m))
				break
			}
		}
	}
	return out
}

// Finding é o resultado de uma mensagem de erro: cada objeto citado e os mods que o trazem.
type Finding struct {
	Object string   `json:"object"`
	Mods   []ModRef `json:"mods"`
}

// Check lê a mensagem e procura os objetos nos mods. SizeMismatch diz se é o erro de versões diferentes.
func Check(text string, mods []Mod) (findings []Finding, sizeMismatch bool) {
	low := strings.ToLower(text)
	sizeMismatch = strings.Contains(low, "serial size mismatch") || strings.Contains(low, "corrupt data found")
	findings = []Finding{}
	for _, obj := range Assets(text) {
		findings = append(findings, Finding{Object: obj, Mods: nonNil(Owners(obj, mods))})
	}
	return findings, sizeMismatch
}

func nonNil(r []ModRef) []ModRef {
	if r == nil {
		return []ModRef{}
	}
	return r
}
