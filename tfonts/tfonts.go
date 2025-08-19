package main

import (
	"flag"
	"image"
	"image/color"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/terminal"
	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	fontDirFlag := flag.String("fontdir", "/System/Library/Fonts", "Directory containing font files")
	cli.Main()
	ap := ansipixels.NewAnsiPixels(60)
	err := ap.Open()
	if err != nil {
		return log.FErrf("failed to open ansi pixels: %v", err)
	}
	defer func() {
		ap.MoveCursor(0, ap.H-1)
		ap.Restore()
	}()
	terminal.LoggerSetup(&terminal.CRLFWriter{Out: ap.Out})
	fdir := *fontDirFlag
	// Walk the font directory - recursively find all fonts and add them to a slice
	var fonts []string
	err = filepath.WalkDir(fdir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Errf("error accessing %s: %v", path, err)
			return err
		}
		if d.IsDir() {
			return nil
		}
		// check for ttf or ttc suffixes
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".ttf" || ext == ".ttc" {
			fonts = append(fonts, path)
		}
		return nil
	})
	log.Infof("Found %d fonts", len(fonts))
	if err != nil {
		return log.FErrf("failed to walk font directory: %v", err)
	}
	for _, f := range fonts {
		col := tcolor.HSLToRGB(rand.Float64(), 0.5, 0.5)
		textColor := color.RGBA{R: col.R, G: col.G, B: col.B, A: 255}
		b, err := os.ReadFile(f)
		if err != nil {
			log.Errf("error reading font %s: %v", f, err)
			continue
		}
		log.Infof("Processing font file: %q", f)
		// Draw the font onto the image
		fc, err := opentype.ParseCollection(b)
		if err != nil {
			log.Errf("failed to parse font %s: %v", f, err)
		} else {
			log.Infof("Loaded font: %s", f)
		}
		var buf sfnt.Buffer
		for i := range fc.NumFonts() {
			face, err := fc.Font(i)
			if err != nil {
				log.Errf("failed to get sub font %s / %d: %v", f, i, err)
				continue
			}
			if i > 0 {
				break // only draw the first font in the collection for now.
			}
			idx, err := face.GlyphIndex(&buf, 'j') // check if the font has basic glyphs
			if err != nil {
				log.Errf("failed to get glyph index for font %s / %d: %v", f, i, err)
				continue
			}
			if idx == 0 {
				log.Infof("Font %s / %d does not have glyph for 'j'", f, i)
				continue
			}
			name, err := face.Name(nil, sfnt.NameIDFull)
			if err != nil {
				log.Errf("failed to get name for font %s / %d: %v", f, i, err)
				continue
			}
			log.Infof("Drawing font %d: %s\n%s", i, f, name)
			offsetY := 6
			offsetX := 3
			ff, err := opentype.NewFace(face, &opentype.FaceOptions{Size: 36, DPI: 72, Hinting: font.HintingFull})
			ap.OnResize = func() error {
				img := image.NewRGBA(image.Rect(0, 0, ap.W, ap.H*2))
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(textColor),
					Face: ff,
					Dot:  fixed.Point26_6{X: fixed.I(offsetX), Y: fixed.I(ap.H - offsetY)},
				}
				d.DrawString("The quick brown fox")
				d.Dot.X = fixed.I(offsetX)
				d.Dot.Y = fixed.I(2*(ap.H-2) - offsetY)
				d.DrawString("jumps over the lazy dog")
				ap.StartSyncMode()
				ap.ClearScreen()
				ap.ShowScaledImage(img)
				ap.WriteAtStr(0, 0, name)
				return nil
			}
			ap.OnResize()
			ap.ReadOrResizeOrSignal()
			if len(ap.Data) > 0 && ap.Data[0] == 'q' {
				log.Infof("Exiting on user request")
				return 0
			}
		}
	}
	return 0
}
