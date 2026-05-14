package main

import (
	"os"
	"runtime"
)

// firstExisting returns the first path that exists on disk, or the last path as fallback.
func firstExisting(paths ...string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return paths[len(paths)-1]
}

// fontPath returns the OS-specific file path for a given font name.
// Supported names: "Impact", "Arial", "Arial Bold".
func fontPath(name string) string {
	switch runtime.GOOS {
	case "windows":
		switch name {
		case "Impact":
			return firstExisting(
				`C:\Windows\Fonts\impact.ttf`,
				`C:\Windows\Fonts\DejaVuSans-Bold.ttf`,
			)
		case "Arial":
			return firstExisting(
				`C:\Windows\Fonts\arial.ttf`,
				`C:\Windows\Fonts\LiberationSans-Regular.ttf`,
			)
		case "Arial Bold":
			return firstExisting(
				`C:\Windows\Fonts\arialbd.ttf`,
				`C:\Windows\Fonts\LiberationSans-Bold.ttf`,
			)
		}
	case "darwin":
		switch name {
		case "Impact":
			return firstExisting(
				"/System/Library/Fonts/Supplemental/Impact.ttf",
				"/System/Library/Fonts/Impact.ttf",
				"/Library/Fonts/Impact.ttf",
			)
		case "Arial":
			return firstExisting(
				"/System/Library/Fonts/Supplemental/Arial.ttf",
				"/System/Library/Fonts/Arial.ttf",
				"/Library/Fonts/Arial.ttf",
			)
		case "Arial Bold":
			return firstExisting(
				"/System/Library/Fonts/Supplemental/Arial Bold.ttf",
				"/System/Library/Fonts/Arial Bold.ttf",
				"/Library/Fonts/Arial Bold.ttf",
			)
		}
	default: // linux, freebsd, etc.
		switch name {
		case "Impact":
			return firstExisting(
				"/usr/share/fonts/truetype/msttcorefonts/Impact.ttf",
				"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
				"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
			)
		case "Arial":
			return firstExisting(
				"/usr/share/fonts/truetype/msttcorefonts/Arial.ttf",
				"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
				"/usr/share/fonts/TTF/LiberationSans-Regular.ttf",
			)
		case "Arial Bold":
			return firstExisting(
				"/usr/share/fonts/truetype/msttcorefonts/Arial_Bold.ttf",
				"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
				"/usr/share/fonts/TTF/LiberationSans-Bold.ttf",
			)
		}
	}
	return name
}
