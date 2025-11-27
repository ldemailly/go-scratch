//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ldemailly/go-scratch/launchclick"
)

func main() {
	var title string
	var rx float64
	var ry float64
	var delayMs int
	var clicks int
	var child string
	var debug bool
	var listTree bool

	flag.StringVar(&title, "title", "", "Window title (full or partial)")
	flag.Float64Var(&rx, "x", 0.5, "Relative x within window (0..1)")
	flag.Float64Var(&ry, "y", 0.5, "Relative y within window (0..1)")
	flag.IntVar(&delayMs, "delay", 300, "Delay ms before clicking after bringing to front")
	flag.IntVar(&clicks, "clicks", 1, "Number of clicks")
	flag.StringVar(&child, "child", "", "Optional: child control title or class to search for inside the found window")
	flag.BoolVar(&debug, "debug", false, "If true and searching children, list discovered child controls")
	flag.BoolVar(&listTree, "list-tree", false, "List full widget tree of the found window")
	flag.Parse()

	if title == "" {
		fmt.Println("-title is required")
		flag.Usage()
		os.Exit(2)
	}

	hwnd, foundTitle := launchclick.FindWindowByTitlePartial(title)
	if hwnd == 0 {
		fmt.Printf("No window found matching %q\n", title)
		os.Exit(1)
	}
	fmt.Printf("Found window: hwnd=0x%X title=%q\n", hwnd, foundTitle)

	if listTree && child == "" {
		tree := launchclick.ListWidgetTree(hwnd)
		fmt.Println("Widget tree:")
		for _, line := range tree {
			fmt.Println(line)
		}
		return
	}

	if child != "" {
		var debugList []string
		// try recursive search
		ch, err := launchclick.FindChildRecursiveByTitle(hwnd, child, debug, &debugList)
		if listTree || debug {
			fmt.Println("Widget tree:")
			tree := launchclick.ListWidgetTree(hwnd)
			for _, line := range tree {
				fmt.Println(line)
			}
		}
		if err != nil {
			fmt.Printf("Child control matching %q not found (recursive): %v\n", child, err)
		} else {
			fmt.Printf("Found child (recursive) hwnd=0x%X\n", ch)
			// bring child to front by bringing parent and then clicking center of child
			launchclick.BringWindowToFront(hwnd)
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			// compute child rect and click center
			rect, err := launchclick.GetWindowRect(ch)
			if err != nil {
				fmt.Printf("GetWindowRect for child failed: %v\n", err)
			} else {
				cx := int(rect.Left + (rect.Right-rect.Left)/2)
				cy := int(rect.Top + (rect.Bottom-rect.Top)/2)
				for i := 0; i < clicks; i++ {
					launchclick.ClickAtAbsolute(cx, cy)
				}
			}
		}
		return
	}

	launchclick.BringWindowToFront(hwnd)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)

	for i := 0; i < clicks; i++ {
		if err := launchclick.ClickAtRelative(hwnd, rx, ry); err != nil {
			fmt.Printf("Click failed: %v\n", err)
			os.Exit(1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
