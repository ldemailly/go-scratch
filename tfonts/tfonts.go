package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
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
	var defaultFontDir string
	switch runtime.GOOS {
	case "darwin":
		defaultFontDir = "/System/Library/Fonts"
	case "linux", "freebsd":
		defaultFontDir = "/usr/share/fonts"
	case "windows":
		defaultFontDir = "C:\\Windows\\Fonts"
	}
	fixedSeed := flag.Int64("seed", 0, "set fixed seed, 0 is random one")
	fontDirFlag := flag.String("fontdir", defaultFontDir, "Directory `path` containing font files")
	fontSizeFlag := flag.Float64("size", 36, "Font size in `points`")
	singleColor := flag.String("color", "", "Single text color, if empty use random colors")
	defaultTruecolor := ansipixels.DetectColorMode().TrueColor
	trueColor := flag.Bool("truecolor", defaultTruecolor, "Use true color (24-bit) instead of 256 colors")
	monoFlag := flag.Bool("mono", false, "Use monochrome (1-bit) color")
	grayFlag := flag.Bool("gray", false, "Use grayscale")
	runeFlag := flag.String("rune", "", "Rune to check for in fonts (default: first `rune` of first line)")
	autoPlayFlag := flag.Duration("autoplay", 0, "If > 0, automatically advance to next font after this duration (e.g. 2s, 500ms)")
	fontFlag := flag.String("font", "", "Font `path` to use instead of showing all the fonts in fontdir")
	allVariantsFlag := flag.Bool("all", false, "Show all font variants (default is only the first found per file)")
	cli.MaxArgs = -1
	cli.ArgsHelp = "2 lines of words to use or default text"
	cli.Main()
	autoPlay := *autoPlayFlag
	var line1, line2 string
	var runeToCheck rune = 0
	switch flag.NArg() {
	case 2:
		line1 = flag.Arg(0)
		line2 = flag.Arg(1)
	case 1:
		input, err := strconv.Unquote(`"` + flag.Arg(0) + `"`)
		if err != nil {
			return log.FErrf("failed to unquote input %q: %v", flag.Arg(0), err)
		}
		lines := strings.SplitN(input, "\n", 2)
		line1 = lines[0]
		line2 = ""
		if len(lines) > 1 {
			line2 = lines[1]
		}
	case 0:
		line1 = "The quick brown fox"
		line2 = "jumps over the lazy dog"
		runeToCheck = 'j' // not T as some symbol fonts have T but not j
	default:
		allInput := strings.Join(flag.Args(), " ")
		mid := len(allInput)/2 - 1
		cutOff := strings.Index(allInput[mid:], " ")
		line1 = allInput[:mid+cutOff]
		line2 = allInput[mid+cutOff+1:]
	}
	if runeToCheck == 0 {
		// use the first rune of line 1
		runeToCheck = []rune(line1)[0]
	}
	if *runeFlag != "" {
		runeToCheck = []rune(*runeFlag)[0]
	}
	fps := 60.
	if autoPlay > 0 {
		fps = 1 / autoPlay.Seconds()
	}
	ap := ansipixels.NewAnsiPixels(fps)
	err := ap.Open()
	if err != nil {
		return log.FErrf("failed to open ansi pixels: %v", err)
	}
	defer func() {
		ap.MoveCursor(0, ap.H-1)
		ap.Restore()
	}()
	ap.TrueColor = *trueColor
	ap.Gray = *grayFlag
	if *monoFlag {
		ap.TrueColor = false
		ap.Color256 = false
		if *singleColor != "" {
			c, err := tcolor.FromString(*singleColor)
			if err != nil {
				return log.FErrf("invalid color %q: %v", *singleColor, err)
			}
			t, v := c.Decode()
			if t != tcolor.ColorTypeBasic {
				return log.FErrf("For mono, use basic color, got %s", c.String())
			}
			ap.MonoColor = tcolor.BasicColor(v)
		}
	}
	terminal.LoggerSetup(&terminal.CRLFWriter{Out: ap.Out})
	ap.SyncBackgroundColor()
	fdir := *fontDirFlag
	// Walk the font directory - recursively find all fonts and add them to a slice
	var fonts []string
	if *fontFlag != "" {
		fonts = append(fonts, *fontFlag)
	} else {
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
	}
	// Set fixed rand/v2 seed:
	// make a reproducible generator with seed 42
	var src rand.Source
	if *fixedSeed > 0 {
		src = rand.NewPCG(uint64(*fixedSeed), 0)
	} else {
		src = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
	rnd := rand.New(src)
	textColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	extra := " " // one space after the cursor
	if !*monoFlag && *singleColor != "" {
		c, err := tcolor.FromString(*singleColor)
		if err != nil {
			return log.FErrf("invalid color %q: %v", *singleColor, err)
		}
		t, v := c.Decode()
		if t == tcolor.ColorTypeBasic || t == tcolor.ColorType256 {
			return log.FErrf("Use HSL or RGB/Hex color, got %s", c.String())
		}
		rgb := tcolor.ToRGB(t, v)
		textColor = color.RGBA{R: rgb.R, G: rgb.G, B: rgb.B, A: 255}
	}
	fidx := 0
	for fidx < len(fonts) {
		if !*monoFlag && *singleColor == "" {
			col := tcolor.HSLToRGB(rnd.Float64(), 0.5, 0.6)
			textColor = color.RGBA{R: col.R, G: col.G, B: col.B, A: 255}
		}
		f := fonts[fidx]
		b, err := os.ReadFile(f)
		if err != nil {
			log.Errf("error reading font %s: %v", f, err)
			fonts = slices.Delete(fonts, fidx, fidx+1)
			continue
		}
		log.LogVf("Processing font file: %q", f)
		// Draw the font onto the image
		fc, err := opentype.ParseCollection(b)
		if err != nil {
			log.Errf("failed to parse font %s: %v", f, err)
			fonts = slices.Delete(fonts, fidx, fidx+1)
			continue
		} else {
			log.LogVf("Loaded font: %s", f)
		}
		var buf sfnt.Buffer
		numSubFonts := fc.NumFonts()
		i := 0
		for i < numSubFonts {
			// TODO refactor the error cases instead of copy pasta
			face, err := fc.Font(i)
			if err != nil {
				log.Errf("failed to get sub font %s / %d: %v", f, i, err)
				fonts = slices.Delete(fonts, fidx, fidx+1)
				fidx--
				break
			}
			if !*allVariantsFlag && i > 0 {
				break // only draw the first font in the collection for now.
			}
			i++
			idx, err := face.GlyphIndex(&buf, runeToCheck) // check if the font has basic glyphs
			if err != nil {
				log.Errf("failed to get glyph index for font %s / %d: %v", f, i, err)
				fonts = slices.Delete(fonts, fidx, fidx+1)
				fidx--
				break
			}
			if idx == 0 {
				log.Infof("Font %s / %d does not have glyph for '%c'", f, i, runeToCheck)
				fonts = slices.Delete(fonts, fidx, fidx+1)
				fidx--
				break
			}
			name, err := face.Name(nil, sfnt.NameIDFull)
			if err != nil {
				log.Errf("failed to get name for font %s / %d: %v", f, i, err)
				fidx--
				break
			}
			log.LogVf("Drawing font %d: %s\n%s", i, f, name)
			offsetY := 6
			offsetX := 3
			ff, err := opentype.NewFace(face, &opentype.FaceOptions{Size: *fontSizeFlag, DPI: 72, Hinting: font.HintingFull})
			ap.OnResize = func() error {
				img := image.NewRGBA(image.Rect(0, 0, ap.W, ap.H*2))
				// fill img with bgcolor using uniform color
				bgColor := color.RGBA{R: ap.Background.R, G: ap.Background.G, B: ap.Background.B, A: 255}
				draw.Draw(img, img.Bounds(), image.NewUniform(bgColor), image.Point{}, draw.Src)
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(textColor),
					Face: ff,
					Dot:  fixed.Point26_6{X: fixed.I(offsetX), Y: fixed.I(ap.H - offsetY)},
				}
				if len(line2) == 0 { // draw line 1 on line 2, allows for taller
					d.Dot.X = fixed.I(offsetX)
					d.Dot.Y = fixed.I(2*(ap.H-2) - offsetY)
				}
				d.DrawString(line1)
				if len(line2) != 0 {
					d.Dot.X = fixed.I(offsetX)
					d.Dot.Y = fixed.I(2*(ap.H-2) - offsetY)
					d.DrawString(line2)
				}
				ap.StartSyncMode()
				ap.ClearScreen()
				ap.ShowScaledImage(img)
				subfontInfo := ""
				if numSubFonts > 1 {
					subfontInfo = fmt.Sprintf("(subfont %d/%d) ", i, numSubFonts)
				}
				ap.WriteAt(0, 0, "%d/%d %s%s%s", fidx+1, len(fonts), subfontInfo, name, extra)
				return nil
			}
			ap.OnResize()
			if autoPlay > 0 {
				ap.EndSyncMode()
				ap.ReadOrResizeOrSignalOnce()
			} else {
				ap.ReadOrResizeOrSignal()
			}
			if len(ap.Data) == 0 {
				continue
			}
			c := ap.Data[0]
			switch c {
			case 'q', 'Q', 3:
				ap.MoveCursor(0, 1)
				log.Infof("Exiting on user request, last font file %s", f)
				return 0
			case 127:
				fidx -= 2       // go back one (will be incremented at end of loop)
				i = numSubFonts // poor man's break
			case 27:
				// left arrow
				if len(ap.Data) >= 3 && ap.Data[2] == 'D' {
					fidx -= 2       // go back one (will be incremented at end of loop)
					i = numSubFonts // poor man's break
				}
			}
		}
		fidx = max(fidx+1, 0)
	}
	if autoPlay > 0 {
		// one last key at the end before exiting
		extra = " (last font, any key to exit)..."
		ap.OnResize()
		ap.ReadOrResizeOrSignal()
	}
	return 0
}
