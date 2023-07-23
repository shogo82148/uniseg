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
	// https://learn.microsoft.com/en-us/windows/win32/intl/code-page-identifiers
	// 932:   ANSI/OEM Japanese; Japanese (Shift-JIS)
	// 51932: EUC Japanese
	// 936:   ANSI/OEM Simplified Chinese (PRC, Singapore); Chinese Simplified (GB2312)
	// 949:   ANSI/OEM Korean (Unified Hangul Code)
	// 950:   ANSI/OEM Traditional Chinese (Taiwan; Hong Kong SAR, PRC); Chinese Traditional (Big5)
	case 932, 51932, 936, 949, 950:
		return true
	}

	return false
}
