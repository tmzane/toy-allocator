package main

import "github.com/fatih/color"

// freeList is an allocator that uses a reserved area of memory for bookkeeping.
type freeList struct {
	memorySize             byte
	freeBlocks, usedBlocks map[byte]byte // addr:size

	lastModifiedRange [2]byte
}

func (fl *freeList) name() string { return "free list" }

func (fl *freeList) init(size byte) {
	fl.memorySize = size
	fl.freeBlocks = map[byte]byte{0: size}
	fl.usedBlocks = map[byte]byte{}
}

func (fl *freeList) malloc(msize byte) {
	for addr, size := range fl.freeBlocks {
		if size < msize {
			continue
		}

		delete(fl.freeBlocks, addr)
		fl.usedBlocks[addr] = msize
		fl.lastModifiedRange = [2]byte{addr, addr + msize}

		if diff := size - msize; diff > 0 {
			fl.freeBlocks[addr+msize] = diff
		}

		return // addr
	}

	panic("out of memory")
}

func (fl *freeList) free(addr byte) {
	size, ok := fl.usedBlocks[addr]
	if !ok {
		panic("invalid memory address")
	}

	delete(fl.usedBlocks, addr)
	fl.lastModifiedRange = [2]byte{addr, addr + size}

	raddr := addr + size
	if rsize, ok := fl.freeBlocks[raddr]; ok {
		size += rsize
		delete(fl.freeBlocks, raddr)
	}

	for laddr, lsize := range fl.freeBlocks {
		if laddr+lsize == addr {
			addr = laddr
			size += lsize
			delete(fl.freeBlocks, laddr)
			break
		}
	}

	fl.freeBlocks[addr] = size
}

func (fl *freeList) printSnapshot() {
	used := make([]bool, fl.memorySize)
	for addr, size := range fl.usedBlocks {
		for i := addr; i < addr+size; i++ {
			used[i] = true
		}
	}

	var cells []string
	for addr := byte(0); addr < fl.memorySize; addr++ {
		var attrs []color.Attribute
		if fl.lastModifiedRange[0] <= addr && addr < fl.lastModifiedRange[1] {
			attrs = append(attrs, color.BlinkSlow)
		}

		if used[addr] {
			cells = append(cells, formatCell("xx", append(attrs, color.FgYellow)...))
		} else {
			cells = append(cells, formatCell("--", attrs...))
		}
	}

	printSnapshot(cells, len(fl.freeBlocks), len(fl.usedBlocks))
}
