//go:build windows
// +build windows

package launchclick

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procEnumChildWindows         = user32.NewProc("EnumChildWindows")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowRect            = user32.NewProc("GetWindowRect")
	procShowWindow               = user32.NewProc("ShowWindow")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procSetCursorPos             = user32.NewProc("SetCursorPos")
	procMouseEvent               = user32.NewProc("mouse_event")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procAttachThreadInput        = user32.NewProc("AttachThreadInput")
	procBringWindowToTop         = user32.NewProc("BringWindowToTop")
	procGetClassNameW            = user32.NewProc("GetClassNameW")
)

const (
	SW_RESTORE           = 9
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP   = 0x0004
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

func getWindowText(hwnd uintptr) string {
	r1, _, _ := procGetWindowTextLengthW.Call(hwnd)
	length := int(r1)
	if length == 0 {
		return ""
	}
	buf := make([]uint16, length+1)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func isWindowVisible(hwnd uintptr) bool {
	r, _, _ := procIsWindowVisible.Call(hwnd)
	return r != 0
}

// FindWindowByTitlePartial returns the first visible window whose title contains query (case-insensitive)
func FindWindowByTitlePartial(query string) (uintptr, string) {
	var found uintptr
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		if !isWindowVisible(hwnd) {
			return 1
		}
		title := getWindowText(hwnd)
		if title == "" {
			return 1
		}
		if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
			found = hwnd
			return 0
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	if found == 0 {
		return 0, ""
	}
	return found, getWindowText(found)
}

func BringWindowToFront(hwnd uintptr) error {
	procShowWindow.Call(hwnd, uintptr(SW_RESTORE))

	var pid uintptr
	tidRaw, _, _ := procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	targetTid := uint32(tidRaw)

	kernel := syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentThreadId := kernel.NewProc("GetCurrentThreadId")
	curTidRaw, _, _ := procGetCurrentThreadId.Call()
	curTid := uint32(curTidRaw)

	procAttachThreadInput.Call(uintptr(curTid), uintptr(targetTid), 1)
	procSetForegroundWindow.Call(hwnd)
	procBringWindowToTop.Call(hwnd)
	procAttachThreadInput.Call(uintptr(curTid), uintptr(targetTid), 0)
	return nil
}

func GetWindowRect(hwnd uintptr) (RECT, error) {
	var r RECT
	ret, _, err := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if ret == 0 {
		return r, err
	}
	return r, nil
}

func ClickAtAbsolute(x, y int) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
	time.Sleep(30 * time.Millisecond)
	procMouseEvent.Call(uintptr(MOUSEEVENTF_LEFTDOWN), 0, 0, 0, 0)
	time.Sleep(30 * time.Millisecond)
	procMouseEvent.Call(uintptr(MOUSEEVENTF_LEFTUP), 0, 0, 0, 0)
}

func ClickAtRelative(hwnd uintptr, rx, ry float64) error {
	rect, err := GetWindowRect(hwnd)
	if err != nil {
		return err
	}
	w := int(rect.Right - rect.Left)
	h := int(rect.Bottom - rect.Top)
	absX := int(rect.Left) + int(float64(w)*rx)
	absY := int(rect.Top) + int(float64(h)*ry)
	ClickAtAbsolute(absX, absY)
	return nil
}

// helper to get class name
func getClassName(hwnd uintptr) string {
	buf := make([]uint16, 256)
	procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

// FindChildByTitle searches child windows of parent hwnd for a control whose window text or class contains query.
// If debug is true, the function collects a slice of discovered child descriptions (hwnd, class, title) in the debugOut pointer.
func FindChildByTitle(parent uintptr, query string, debug bool, debugOut *[]string) (uintptr, error) {
	var found uintptr
	lowerQ := strings.ToLower(query)
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		title := getWindowText(hwnd)
		class := getClassName(hwnd)
		if debug && debugOut != nil {
			*debugOut = append(*debugOut, formatChild(hwnd, class, title))
		}
		if strings.Contains(strings.ToLower(title), lowerQ) || strings.Contains(strings.ToLower(class), lowerQ) {
			found = hwnd
			return 0
		}
		return 1
	})
	procEnumChildWindows.Call(parent, cb, 0)
	if found == 0 {
		return 0, errors.New("child not found")
	}
	return found, nil
}

func formatChild(hwnd uintptr, class, title string) string {
	return fmt.Sprintf("hwnd=0x%X class=%q title=%q", hwnd, class, title)
}

// recursive helper: enumerate children and call fn(hwnd, depth) for each
func enumerateChildrenRecursive(parent uintptr, depth int, fn func(hwnd uintptr, depth int) bool) {
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		cont := fn(hwnd, depth)
		if !cont {
			return 0
		}
		// recurse into this child's children
		enumerateChildrenRecursive(hwnd, depth+1, fn)
		return 1
	})
	procEnumChildWindows.Call(parent, cb, 0)
}

// ListWidgetTree returns a slice of strings describing the widget tree rooted at parent
func ListWidgetTree(parent uintptr) []string {
	var out []string
	// include the root
	out = append(out, formatChild(parent, getClassName(parent), getWindowText(parent)))
	enumerateChildrenRecursive(parent, 1, func(hwnd uintptr, depth int) bool {
		indent := strings.Repeat("  ", depth)
		out = append(out, indent+formatChild(hwnd, getClassName(hwnd), getWindowText(hwnd)))
		return true
	})
	return out
}

// FindChildRecursiveByTitle searches recursively for a descendant control whose title or class contains query.
// If debugOut != nil, it will append visited nodes (formatted) to debugOut.
func FindChildRecursiveByTitle(parent uintptr, query string, debug bool, debugOut *[]string) (uintptr, error) {
	var found uintptr
	lowerQ := strings.ToLower(query)
	enumerateChildrenRecursive(parent, 1, func(hwnd uintptr, depth int) bool {
		title := getWindowText(hwnd)
		class := getClassName(hwnd)
		if debug && debugOut != nil {
			indent := strings.Repeat("  ", depth)
			*debugOut = append(*debugOut, indent+formatChild(hwnd, class, title))
		}
		if strings.Contains(strings.ToLower(title), lowerQ) || strings.Contains(strings.ToLower(class), lowerQ) {
			found = hwnd
			return false // stop
		}
		return true
	})
	if found == 0 {
		return 0, errors.New("child not found")
	}
	return found, nil
}
