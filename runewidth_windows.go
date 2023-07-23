//go:build windows

// comes from https://github.com/mattn/go-runewidth/blob/2c6a438f68cfe01255a90824599da41fdf76d1e2/runewidth_windows.go

package uniseg

import (
	"syscall"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32")
	procGetConsoleOutputCP = kernel32.NewProc("GetConsoleOutputCP")
)

// IsEastAsian return true if the current locale is CJK
func IsEastAsian() bool {
	r1, _, _ := procGetConsoleOutputCP.Call()
	if r1 == 0 {
		return false
	}

	switch int(r1) {
	case 932, 51932, 936, 949, 950:
		return true
	}

	return false
}
