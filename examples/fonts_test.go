package main

import "runtime"

// fontPath returns the OS-specific file path for a given font name.
// Supported names: "Impact", "Arial", "Arial Bold".
func fontPath(name string) string {
	switch runtime.GOOS {
	case "windows":
		switch name {
		case "Impact":
			return `C:\Windows\Fonts\impact.ttf`
		case "Arial":
			return `C:\Windows\Fonts\arial.ttf`
		case "Arial Bold":
			return `C:\Windows\Fonts\arialbd.ttf`
		}
	case "darwin":
		switch name {
		case "Impact":
			return "/System/Library/Fonts/Supplemental/Impact.ttf"
		case "Arial":
			return "/System/Library/Fonts/Supplemental/Arial.ttf"
		case "Arial Bold":
			return "/System/Library/Fonts/Supplemental/Arial Bold.ttf"
		}
	default: // linux, freebsd, etc.
		switch name {
		case "Impact":
			return "/usr/share/fonts/truetype/msttcorefonts/Impact.ttf"
		case "Arial":
			return "/usr/share/fonts/truetype/msttcorefonts/Arial.ttf"
		case "Arial Bold":
			return "/usr/share/fonts/truetype/msttcorefonts/Arial_Bold.ttf"
		}
	}
	return name
}
