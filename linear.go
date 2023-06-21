package main

import "github.com/fatih/color"

// linear is the simplest allocator that does not support free.
type linear struct {
	memorySize byte
	prevAddr   byte
	nextAddr   byte
	totalUsed  int
}

func (l *linear) name() string { return "linear" }

func (l *linear) init(size byte) { l.memorySize = size }

func (l *linear) malloc(size byte) {
	if l.memorySize-l.nextAddr < size {
		panic("out of memory")
	}

	l.totalUsed++
	l.prevAddr = l.nextAddr
	l.nextAddr += size
	//	return l.prevAddr
}

func (l *linear) free(addr byte) { panic("not supported") }

func (l *linear) printSnapshot() {
	var cells []string

	for i := byte(0); i < l.nextAddr; i++ {
		if i >= l.prevAddr {
			cells = append(cells, formatCell("xx", color.FgYellow, color.BlinkSlow))
		} else {
			cells = append(cells, formatCell("xx", color.FgYellow))
		}
	}

	for i := l.nextAddr; i < l.memorySize; i++ {
		cells = append(cells, formatCell("--"))
	}

	totalFree := 1
	if l.nextAddr == l.memorySize {
		totalFree = 0
	}

	printSnapshot(cells, totalFree, l.totalUsed)
}
