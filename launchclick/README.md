LaunchClick (Windows)
======================

Build
-----

Requires Go 1.20+ (module uses go 1.24). Build on Windows.

Build the CLI specifically:

```powershell
go build ./cmd/launchclick
```

This produces `launchclick.exe` in the current directory (because the package directory is `cmd/launchclick`).

Usage
-----

```powershell
.\launchclick.exe -title "notepad" -x 0.5 -y 0.2 -delay 300 -clicks 1
```

Flags:
- `-title` (required) : window title full or partial, case-insensitive
- `-x`, `-y` : relative coordinates inside the window (0..1)
- `-delay` : milliseconds to wait after bringing the window to front before clicking
- `-clicks` : number of clicks to perform

Notes and limitations
---------------------
- Works only on Windows (uses Win32 APIs). File is build-tagged for windows.
- Uses a simple string match on visible windows. If multiple windows match, the first found is used.
- Click is simulated by moving the cursor and emitting mouse events. This will affect the global cursor.
- If SetForegroundWindow fails due to focus rules, the program uses AttachThreadInput trick but behavior may still be limited by OS focus privilege.
