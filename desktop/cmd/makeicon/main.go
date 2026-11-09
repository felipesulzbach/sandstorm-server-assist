// makeicon desenha o ícone do app (escudo areia com mira, o mesmo desenho do AppIcon.tsx) e grava
// build/appicon.png (1024 px) e build/windows/icon.ico (16 a 256 px, PNG dentro do ICO).
// Uso: go run ./cmd/makeicon (em desktop/). Base: cmd/makeicon do midia-multi-tools.
package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

func main() {
	must(os.MkdirAll("build", 0o750))
	must(writePNG(filepath.Join("build", "appicon.png"), render(1024)))

	var pngs [][]byte
	sizes := []int{16, 24, 32, 48, 64, 128, 256}
	for _, s := range sizes {
		var buf bytes.Buffer
		must(png.Encode(&buf, render(s)))
		pngs = append(pngs, buf.Bytes())
	}
	must(os.MkdirAll(filepath.Join("build", "windows"), 0o750))
	must(os.WriteFile(filepath.Join("build", "windows", "icon.ico"), ico(sizes, pngs), 0o600))
	log.Println("ícone gerado: build/appicon.png e build/windows/icon.ico")
}

// render desenha o ícone em size×size com 4×4 amostras por pixel (bordas suaves).
func render(size int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	const ss = 4
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var r, g, b, a float64
			for sy := 0; sy < ss; sy++ {
				for sx := 0; sx < ss; sx++ {
					u := (float64(x) + (float64(sx)+0.5)/ss) / float64(size)
					v := (float64(y) + (float64(sy)+0.5)/ss) / float64(size)
					c := shade(u, v)
					r += c[0] * c[3]
					g += c[1] * c[3]
					b += c[2] * c[3]
					a += c[3]
				}
			}
			n := float64(ss * ss)
			if a > 0 {
				r, g, b = r/a, g/a, b/a
			}
			img.SetNRGBA(x, y, color.NRGBA{R: u8(r), G: u8(g), B: u8(b), A: u8(a / n)})
		}
	}
	return img
}

var (
	sand = [4]float64{0xe0 / 255.0, 0xa4 / 255.0, 0x3a / 255.0, 1}
	dark = [4]float64{0x1a / 255.0, 0x14 / 255.0, 0x08 / 255.0, 1}
)

// shade devolve RGBA (0–1) no ponto (u, v). O desenho usa as coordenadas 0–24 do SVG do AppIcon.
func shade(u, v float64) [4]float64 {
	// Fundo: quadrado arredondado escuro, com degradê leve.
	if roundedRect(u, v, 0.04, 0.04, 0.96, 0.96, 0.2) > 0 {
		return [4]float64{}
	}
	t := (u + v) / 2
	bg := [4]float64{lerp(0.15, 0.10, t), lerp(0.14, 0.09, t), lerp(0.12, 0.08, t), 1}

	// Mapeia o quadrado útil (0.12–0.88) para o espaço 0–24 do SVG.
	x := (u - 0.12) / 0.76 * 24
	y := (v - 0.12) / 0.76 * 24
	if !inShield(x, y) {
		return bg
	}
	// Mira: anel de raio 4.2 (espessura 1.6) e quatro traços.
	d := math.Hypot(x-12, y-11.5)
	if math.Abs(d-4.2) <= 0.8 {
		return dark
	}
	if (math.Abs(x-12) <= 0.8 && ((y >= 5.8 && y <= 8.8) || (y >= 14.2 && y <= 17.2))) ||
		(math.Abs(y-11.5) <= 0.8 && ((x >= 6.3 && x <= 9.3) || (x >= 14.7 && x <= 17.7))) {
		return dark
	}
	return sand
}

// inShield: escudo do AppIcon (M12 2 4 5v6c0 5 3.4 9.3 8 11 4.6-1.7 8-6 8-11V5l-8-3z), aproximado.
func inShield(x, y float64) bool {
	switch {
	case y < 2 || y > 22:
		return false
	case y < 5: // topo inclinado
		return math.Abs(x-12) <= (y-2)/3*8
	case y <= 11: // laterais retas
		return x >= 4 && x <= 20
	default: // parte de baixo arredondada até a ponta
		k := (y - 11) / 11
		return math.Abs(x-12) <= 8*math.Sqrt(math.Max(0, 1-k*k))
	}
}

// roundedRect: distância com sinal (≤ 0 = dentro) até um retângulo de cantos arredondados.
func roundedRect(u, v, x0, y0, x1, y1, r float64) float64 {
	cx, cy := (x0+x1)/2, (y0+y1)/2
	hx, hy := (x1-x0)/2-r, (y1-y0)/2-r
	dx := math.Max(math.Abs(u-cx)-hx, 0)
	dy := math.Max(math.Abs(v-cy)-hy, 0)
	return math.Hypot(dx, dy) - r
}

// ico monta um arquivo .ico com imagens PNG (aceito desde o Windows Vista).
func ico(sizes []int, pngs [][]byte) []byte {
	var buf bytes.Buffer
	le := binary.LittleEndian
	if len(pngs) > 255 {
		log.Fatal("imagens demais para um ICO")
	}
	_ = binary.Write(&buf, le, [3]uint16{0, 1, uint16(len(pngs))}) // #nosec G115 -- len(pngs) <= 255 conferido acima
	offset := 6 + 16*len(pngs)
	for i, data := range pngs {
		var dim byte // 0 = 256 no formato ICO
		if s := sizes[i]; s > 0 && s < 256 {
			dim = byte(s)
		}
		_ = binary.Write(&buf, le, struct {
			W, H, Colors, Reserved byte
			Planes, BPP            uint16
			Size, Offset           uint32
		}{dim, dim, 0, 0, 1, 32, uint32(len(data)), uint32(offset)}) // #nosec G115 -- tamanhos pequenos e fixos
		offset += len(data)
	}
	for _, data := range pngs {
		buf.Write(data)
	}
	return buf.Bytes()
}

func writePNG(path string, img image.Image) error {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o600)
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func u8(f float64) uint8 { return uint8(math.Round(math.Max(0, math.Min(1, f)) * 255)) }

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
