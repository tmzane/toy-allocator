package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func formatCell[T byte | string](value T, attrs ...color.Attribute) string {
	return "[" + color.New(attrs...).Sprintf("%.2v", value) + "]"
}

func printSnapshot(cells []string, totalFree, totalUsed int) {
	// first, calculate the optimal number of cell rows based on the user's terminal width.

	haveWidth := terminalWidth()
	needWidth := 4*len(cells) + len(cells) - 1 // 4 chars per cell * N cells + (N-1) separators.

	rowSize := len(cells)
	for i := 1; i <= len(cells); i *= 2 {
		if needWidth/i <= haveWidth {
			rowSize = len(cells) / i
			break
		}
	}

	header := make([]string, len(cells))
	for i := range cells {
		header[i] = fmt.Sprintf("%#.2x", i)
	}

	for i := 0; i < len(cells); i += rowSize {
		fmt.Println(strings.Join(header[i:i+rowSize], " "))
		fmt.Println(strings.Join(cells[i:i+rowSize], " "))
	}

	fmt.Printf("---\nTotal blocks: %s free; %s used\n",
		color.GreenString("%d", totalFree),
		color.RedString("%d", totalUsed),
	)
}

func terminalWidth() int {
	fd := int(os.Stdout.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil {
		return 80
	}
	return width
}
