package cmdpkg

import (
	"runtime"
	"strings"
)

type __colors struct {
	RED    string // 31
	GREEN  string // 32
	YELLOW string // 33
	BLUE   string // 34
	PINK   string // 35

	BOLD string // 1

	RESET string // 0
}

var (
	COLORS __colors
)

func init() {
	if runtime.GOOS == "linux" {
		COLORS.RED = "\033[31m"
		COLORS.GREEN = "\033[32m"
		COLORS.YELLOW = "\033[33m"
		COLORS.BLUE = "\033[34m"
		COLORS.PINK = "\033[35m"
		COLORS.BOLD = "\033[1m"
		COLORS.RESET = "\033[0m"
	}
}

func ParseColors(s string) string {
	r := strings.NewReplacer(
		".RED", COLORS.RED,
		".GREEN", COLORS.GREEN,
		".YELLOW", COLORS.YELLOW,
		".BLUE", COLORS.BLUE,
		".PINK", COLORS.PINK,
		".BOLD", COLORS.BOLD,
		".RESET", COLORS.RESET,
	)
	return r.Replace(s)
}
