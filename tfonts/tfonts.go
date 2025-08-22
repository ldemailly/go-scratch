package main

import (
	"errors"
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
	"time"

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

type FState struct {
	// Input
	FontFile    string
	FontDir     string
	FontSize    float64
	Line1       string
	Line2       string
	RuneToCheck rune
	AutoPlay    time.Duration
	AllVariants bool
	FixedSeed   int64
	Monochrome  bool
	SingleColor string
	// Internal
	rnd       *rand.Rand
	ap        *ansipixels.AnsiPixels
	fonts     []string
	fidx      int
	textColor color.RGBA
	done      bool
	goBack    bool
	total     int
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
	fs := FState{
		FontFile:    *fontFlag,
		FontDir:     *fontDirFlag,
		FontSize:    *fontSizeFlag,
		AutoPlay:    *autoPlayFlag,
		AllVariants: *allVariantsFlag,
		FixedSeed:   *fixedSeed,
		Monochrome:  *monoFlag,
		SingleColor: *singleColor,
	}
	err := fs.LinesAndRune(*runeFlag)
	if err != nil {
		return log.FErrf("failed to get lines and rune: %v", err)
	}
	fps := 60.
	if fs.AutoPlay > 0 {
		fps = 1 / fs.AutoPlay.Seconds()
	}
	fs.ap = ansipixels.NewAnsiPixels(fps)
	err = fs.ap.Open()
	if err != nil {
		return log.FErrf("failed to open ansi pixels: %v", err)
	}
	defer func() {
		fs.ap.MoveCursor(0, fs.ap.H-1)
		fs.ap.Restore()
	}()
	fs.ap.TrueColor = *trueColor
	fs.ap.Gray = *grayFlag
	if *monoFlag {
		fs.ap.TrueColor = false
		fs.ap.Color256 = false
		if *singleColor != "" {
			c, err := tcolor.FromString(*singleColor)
			if err != nil {
				return log.FErrf("invalid color %q: %v", *singleColor, err)
			}
			t, v := c.Decode()
			if t != tcolor.ColorTypeBasic {
				return log.FErrf("For mono, use basic color, got %s", c.String())
			}
			fs.ap.MonoColor = tcolor.BasicColor(v)
		}
	}
	terminal.LoggerSetup(&terminal.CRLFWriter{Out: fs.ap.Out})
	fs.ap.SyncBackgroundColor()
	if err = fs.FontsList(); err != nil {
		return log.FErrf("failed to list fonts: %v", err)
	}
	log.Infof("Found %d fonts", len(fs.fonts))
	if err = fs.ColorSetup(); err != nil {
		return log.FErrf("failed to set up color(s): %v", err)
	}
	return fs.RunFonts()
}

func (fs *FState) LinesAndRune(runeFlag string) error {
	switch flag.NArg() {
	case 2:
		fs.Line1 = flag.Arg(0)
		fs.Line2 = flag.Arg(1)
	case 1:
		input, err := strconv.Unquote(`"` + flag.Arg(0) + `"`)
		if err != nil {
			return fmt.Errorf("failed to unquote input %q: %w", flag.Arg(0), err)
		}
		lines := strings.SplitN(input, "\n", 2)
		fs.Line1 = lines[0]
		fs.Line2 = ""
		if len(lines) > 1 {
			fs.Line2 = lines[1]
		}
	case 0:
		fs.Line1 = "The quick brown fox"
		fs.Line2 = "jumps over the lazy dog"
		fs.RuneToCheck = 'j' // not T as some symbol fonts have T but not j
	default:
		allInput := strings.Join(flag.Args(), " ")
		mid := len(allInput)/2 - 1
		cutOff := strings.Index(allInput[mid:], " ")
		fs.Line1 = allInput[:mid+cutOff]
		fs.Line2 = allInput[mid+cutOff+1:]
	}
	if fs.RuneToCheck == 0 {
		// use the first rune of line 1
		fs.RuneToCheck = []rune(fs.Line1)[0]
	}
	if runeFlag != "" {
		fs.RuneToCheck = []rune(runeFlag)[0]
	}
	return nil
}

func (fs *FState) FontsList() error {
	// Walk the font directory - recursively find all fonts and add them to a slice
	if fs.FontFile != "" {
		fs.fonts = append(fs.fonts, fs.FontFile)
		return nil
	}
	return filepath.WalkDir(fs.FontDir, func(path string, d os.DirEntry, err error) error {
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
			fs.fonts = append(fs.fonts, path)
		}
		return nil
	})
}

func (fs *FState) ColorSetup() error {
	// Set fixed rand/v2 seed:
	// make a reproducible generator with seed 42
	var src rand.Source
	if fs.FixedSeed > 0 {
		src = rand.NewPCG(uint64(fs.FixedSeed), 0)
	} else {
		src = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}
	fs.rnd = rand.New(src)
	fs.textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if !fs.Monochrome && fs.SingleColor != "" {
		c, err := tcolor.FromString(fs.SingleColor)
		if err != nil {
			return fmt.Errorf("invalid color %q: %w", fs.SingleColor, err)
		}
		t, v := c.Decode()
		if t == tcolor.ColorTypeBasic || t == tcolor.ColorType256 {
			return fmt.Errorf("use HSL or RGB/Hex color, got %s", c.String())
		}
		rgb := tcolor.ToRGB(t, v)
		fs.textColor = color.RGBA{R: rgb.R, G: rgb.G, B: rgb.B, A: 255}
	}
	return nil
}

func (fs *FState) ErrorDeleteEntry(err error, fmt string, args ...any) {
	if errors.Is(err, NoGlyphErr) {
		log.Infof(fmt, args...)
	} else {
		log.Errf(fmt, args...)
	}
	fs.fonts = slices.Delete(fs.fonts, fs.fidx, fs.fidx+1)
}

func (fs *FState) RunFonts() int {
	for fs.fidx < len(fs.fonts) {
		// optionally change color each time
		if !fs.Monochrome && fs.SingleColor == "" {
			col := tcolor.HSLToRGB(fs.rnd.Float64(), 0.5, 0.6)
			fs.textColor = color.RGBA{R: col.R, G: col.G, B: col.B, A: 255}
		}
		// Do the work for 1 font
		if err := fs.ProcessOneFile(); err != nil {
			fs.ErrorDeleteEntry(err, "error processing font %s: %v", fs.fonts[fs.fidx], err)
			continue
		}
		if fs.done {
			return 0
		}
		if fs.goBack {
			fs.fidx = max(fs.fidx-1, 0)
		} else {
			fs.fidx++
		}
	}
	fs.ap.MoveCursor(0, 1)
	fs.ap.EndSyncMode()
	log.Infof("Total fonts processed: %d", fs.total)
	return 0
}

func (fs *FState) ProcessOneFile() error {
	fs.goBack = false
	f := fs.fonts[fs.fidx]
	b, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	log.LogVf("Processing font file: %q", f)
	// Draw the font onto the image
	fc, err := opentype.ParseCollection(b)
	if err != nil {
		return err
	}
	log.LogVf("Loaded font: %s", f)
	numSubFonts := fc.NumFonts()
	i := 0
	for i < numSubFonts {
		if err := fs.ProcessSubFont(fc, i, numSubFonts); err != nil {
			return err
		}
		i++
		if !fs.AllVariants {
			i = numSubFonts // poor mans 'break after this one loop but process keys'
		}
		if len(fs.ap.Data) == 0 {
			continue
		}
		c := fs.ap.Data[0]
		switch c {
		case 'q', 'Q', 3:
			fs.ap.MoveCursor(0, 1)
			log.Infof("Exiting on user request after %d %s. Last font file %s", fs.total, cli.Plural(fs.total, "font"), f)
			fs.done = true
			return nil
		case 127:
			fs.goBack = true
			return nil
		case 27:
			// left arrow
			if len(fs.ap.Data) >= 3 && fs.ap.Data[2] == 'D' {
				fs.goBack = true
				return nil
			}
		}
	}
	return nil
}

var NoGlyphErr = errors.New("no glyph for rune")

func (fs *FState) ProcessSubFont(fc *opentype.Collection, i, numSubFonts int) error {
	var buf sfnt.Buffer
	face, err := fc.Font(i) // 0 indexed here, we use i+1 for user messages/logs below.
	if err != nil {
		return err
	}
	idx, err := face.GlyphIndex(&buf, fs.RuneToCheck) // check if the font has basic glyphs
	if err != nil {
		return err
	}
	if idx == 0 {
		return fmt.Errorf("%w %c in font %s / %d", NoGlyphErr, fs.RuneToCheck, fs.fonts[fs.fidx], i+1)
	}
	name, err := face.Name(nil, sfnt.NameIDFull)
	if err != nil {
		return err
	}
	log.LogVf("Drawing font %d: %s\n%s", i+1, fs.fonts[fs.fidx], name)
	offsetY := 6
	offsetX := 3
	ff, err := opentype.NewFace(face, &opentype.FaceOptions{Size: fs.FontSize, DPI: 72, Hinting: font.HintingFull})
	fs.ap.OnResize = func() error {
		img := image.NewRGBA(image.Rect(0, 0, fs.ap.W, fs.ap.H*2))
		// fill img with bgcolor using uniform color
		bgColor := color.RGBA{R: fs.ap.Background.R, G: fs.ap.Background.G, B: fs.ap.Background.B, A: 255}
		draw.Draw(img, img.Bounds(), image.NewUniform(bgColor), image.Point{}, draw.Src)
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(fs.textColor),
			Face: ff,
			Dot:  fixed.Point26_6{X: fixed.I(offsetX), Y: fixed.I(fs.ap.H - offsetY)},
		}
		if len(fs.Line2) == 0 { // draw line 1 on line 2, allows for taller
			d.Dot.X = fixed.I(offsetX)
			d.Dot.Y = fixed.I(2*(fs.ap.H-2) - offsetY)
		}
		d.DrawString(fs.Line1)
		if len(fs.Line2) != 0 {
			d.Dot.X = fixed.I(offsetX)
			d.Dot.Y = fixed.I(2*(fs.ap.H-2) - offsetY)
			d.DrawString(fs.Line2)
		}
		fs.ap.StartSyncMode()
		fs.ap.ClearScreen()
		fs.ap.ShowScaledImage(img)
		subfontInfo := ""
		if numSubFonts > 1 {
			subfontInfo = fmt.Sprintf("(subfont %d/%d) ", i+1, numSubFonts)
		}
		fs.ap.WriteAt(0, 0, "%d/%d %s%s ", fs.fidx+1, len(fs.fonts), subfontInfo, name)
		return nil
	}
	fs.total++
	fs.ap.OnResize()
	if fs.AutoPlay > 0 {
		fs.ap.EndSyncMode()
		fs.ap.ReadOrResizeOrSignalOnce()
	} else {
		fs.ap.ReadOrResizeOrSignal()
	}
	return nil
}
