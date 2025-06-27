package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

/* ⬇️  Configurações básicas */
const (
	inputYAML = "structure.yaml"
	outputDir = "output"
	txtOut    = "output/folder_structure.txt"
	imgOut    = "output/folder_structure.png"
)

/* ⬇️  Temas (dark, light, dracula) */
type Theme struct{ Bg, Fg color.RGBA }

func rgba(r, g, b uint8) color.RGBA { return color.RGBA{r, g, b, 255} }

var themes = map[string]Theme{
	"dark":    {rgba(30, 30, 30), rgba(0, 255, 128)},
	"light":   {rgba(255, 255, 255), rgba(0, 0, 0)},
	"dracula": {rgba(40, 42, 54), rgba(80, 250, 123)},
}

/* ⬇️  Entrada do programa */
func main() {
	themeFlag := flag.String("theme", "dark", "Tema: dark | light | dracula")
	flag.Parse()

	theme, ok := themes[*themeFlag]
	if !ok {
		log.Fatalf("Tema inválido: %s", *themeFlag)
	}

	// Ler YAML ---------------------------------------------------------
	raw, err := ioutil.ReadFile(inputYAML)
	if err != nil {
		log.Fatalf("Não foi possível ler %s: %v", inputYAML, err)
	}

	var tree map[string]interface{}
	if err := yaml.Unmarshal(raw, &tree); err != nil {
		log.Fatalf("YAML inválido: %v", err)
	}

	// Construir linhas -------------------------------------------------
	var lines []string
	buildTree(tree, "", &lines)

	// Salvar TXT -------------------------------------------------------
	_ = os.MkdirAll(outputDir, 0755)
	_ = ioutil.WriteFile(txtOut, []byte(strings.Join(lines, "\n")), 0644)

	// Salvar PNG -------------------------------------------------------
	if err := renderPNG(lines, imgOut, theme); err != nil {
		log.Fatalf("Erro gerando PNG: %v", err)
	}

	fmt.Println("✅ Arquivos gerados em", outputDir)
}

/* ---------- Construção da árvore (pastas antes de arquivos) ---------- */
func buildTree(node interface{}, prefix string, lines *[]string) {
	switch n := node.(type) {

	case map[string]interface{}:
		// Arquivos via __files__
		if rawFiles, ok := n["__files__"].([]interface{}); ok {
			for i, v := range rawFiles {
				if fn, ok := v.(string); ok {
					*lines = append(*lines, prefix+branch(i == len(rawFiles)-1)+" "+fn)
				}
			}
		}

		// Separar subpastas e arquivos-nulos
		var dirs, files []string
		for k, v := range n {
			if k == "__files__" {
				continue
			}
			if v == nil {
				files = append(files, k)
			} else {
				dirs = append(dirs, k)
			}
		}
		sort.Strings(dirs)
		sort.Strings(files)

		total := len(dirs) + len(files)

		// Subpastas primeiro
		for i, d := range dirs {
			last := i == total-1 && len(files) == 0
			*lines = append(*lines, prefix+branch(last)+d+"/")
			buildTree(n[d], prefix+nextPrefix(last), lines)
		}

		// Depois arquivos isolados
		for j, f := range files {
			*lines = append(*lines, prefix+branch(j == len(files)-1)+f)
		}

	case []interface{}: // lista pura de arquivos
		for i, v := range n {
			if fn, ok := v.(string); ok {
				*lines = append(*lines, prefix+branch(i == len(n)-1)+" "+fn)
			}
		}
	}
}

/* ---------- Helpers visuais ---------- */
func branch(last bool) string {
	if last {
		return "└─"
	}
	return "├─"
}
func nextPrefix(last bool) string {
	if last {
		return "   "
	}
	return "│  "
}

/* ---------- Renderização da imagem PNG ---------- */
func renderPNG(lines []string, outPath string, th Theme) error {
	face := basicfont.Face7x13
	lineH := face.Metrics().Height.Ceil() + 6
	maxW := 0
	for _, l := range lines {
		w := font.MeasureString(face, l).Ceil()
		if w > maxW {
			maxW = w
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, maxW+40, lineH*len(lines)+40))
	draw.Draw(img, img.Bounds(), &image.Uniform{th.Bg}, image.Point{}, draw.Src)

	d := &font.Drawer{Dst: img, Src: &image.Uniform{th.Fg}, Face: face}
	y := 20 + face.Metrics().Ascent.Ceil()
	for _, l := range lines {
		d.Dot = fixed.P(20, y)
		d.DrawString(l)
		y += lineH
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
