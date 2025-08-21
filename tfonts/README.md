# tfonts

Demonstrates large true type font rending in the terminal with multiple options

From truecolor to mono

Using [ansipixels](https://github.com/fortio/terminal#fortioorgterminalansipixels) image to terminal rendering.

## Install

```sh
go install github.com/ldemailly/go-scratch/tfonts@latest
```

## Example

```
tfonts -autoplay 0.3s Abc def
```

![Example Screenshot](screenshot.png)


## Options

```
tfonts help
```

```
tfonts v1.0.0 usage:
        tfonts [flags] 2 lines of words to use or default text
or 1 of the special arguments
        tfonts {help|envhelp|version|buildinfo}
flags:
  -autoplay duration
        If > 0, automatically advance to next font after this duration (e.g. 2s, 500ms)
  -color string
        Single text color, if empty use random colors
  -fontdir string
        Directory containing font files (default "/System/Library/Fonts")
  -gray
        Use grayscale
  -mono
        Use monochrome (1-bit) color
  -quiet
        Quiet mode, sets loglevel to Error (quietly) to reduces the output
  -rune string
        Rune to check for in fonts (default: first rune of first line)
  -seed int
        set fixed seed, 0 is random one
  -size float
        Font size in points (default 36)
  -truecolor
        Use true color (24-bit) instead of 256 colors (default true)
```
