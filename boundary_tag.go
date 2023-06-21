package main

import "github.com/fatih/color"

// memory block structure
//
// addr + 0: the size of the block
// addr + 1: is the block free (1) or used (2)
// addr + 2: the beginning of usable memory
// addr + 2 + size: the size of the block again

const (
	flagFree byte = 1
	flagUsed byte = 2

	metadataSize byte = 3 // total bytes for bookkeeping.
)

type block struct {
	addr byte
	size byte
	flag byte
}

// boundaryTag is an allocator that uses the allocated memory itself for bookkeeping.
type boundaryTag struct {
	memory []byte

	lastModifiedRange [2]byte
}

func (bt *boundaryTag) name() string { return "boundary tag" }

func (bt *boundaryTag) init(size byte) {
	bt.memory = make([]byte, size)
	bt.writeBlock(block{
		addr: 0,
		size: size - metadataSize,
		flag: flagFree,
	})
}

func (bt *boundaryTag) malloc(size byte) {
	for addr := byte(0); addr < byte(len(bt.memory)); {
		b := bt.readBlock(addr)
		if b.flag == flagUsed || b.size < size+metadataSize {
			addr += b.size + metadataSize
			continue
		}

		bt.writeBlock(block{addr, size, flagUsed})
		bt.lastModifiedRange = [2]byte{addr + 2, addr + 2 + size}

		if diff := b.size - size; diff > 0 {
			bt.writeBlock(block{
				addr: addr + size + metadataSize,
				size: diff - metadataSize,
				flag: flagFree,
			})
		}

		return // addr + 2
	}

	panic("out of memory")
}

func (bt *boundaryTag) free(addr byte) {
	// initially, `addr` here is the beginning of usable memory.
	addr -= 2

	if bt.memory[addr+1] != flagUsed {
		panic("invalid memory address")
	}

	b := bt.readBlock(addr)
	bt.lastModifiedRange = [2]byte{b.addr + 2, b.addr + 2 + b.size}

	if isLast := b.addr+b.size+metadataSize == byte(len(bt.memory)); !isLast {
		rb := bt.readBlock(b.addr + b.size + metadataSize)
		if rb.flag == flagFree {
			b.size += rb.size + metadataSize
		}
	}

	if isFirst := b.addr == 0; !isFirst {
		lsize := bt.memory[b.addr-1]
		lb := bt.readBlock(b.addr - lsize - metadataSize)
		if lb.flag == flagFree {
			b.addr = lb.addr
			b.size += lb.size + metadataSize
		}
	}

	bt.writeBlock(block{
		addr: b.addr,
		size: b.size,
		flag: flagFree,
	})
}

func (bt *boundaryTag) printSnapshot() {
	var cells []string
	var totalFree, totalUsed int

	for addr := byte(0); addr < byte(len(bt.memory)); {
		b := bt.readBlock(addr)

		// 1. size
		cells = append(cells, formatCell(b.size, color.FgBlue))

		var cellValue string
		var cellAttrs []color.Attribute

		// 2. flag
		if b.flag == flagFree {
			totalFree++
			cellValue = "--"
			cells = append(cells, formatCell(flagFree, color.FgGreen))
		} else {
			totalUsed++
			cellValue = "xx"
			cellAttrs = append(cellAttrs, color.FgYellow)
			cells = append(cells, formatCell(flagUsed, color.FgRed))
		}

		// 3. usable memory
		for i := b.addr + 2; i < b.addr+2+b.size; i++ {
			if bt.lastModifiedRange[0] <= i && i < bt.lastModifiedRange[1] {
				cells = append(cells, formatCell(cellValue, append(cellAttrs, color.BlinkSlow)...))
			} else {
				cells = append(cells, formatCell(cellValue, cellAttrs...))
			}
		}

		// 4. size again
		cells = append(cells, formatCell(b.size, color.FgBlue))

		addr += b.size + metadataSize
	}

	printSnapshot(cells, totalFree, totalUsed)
}

func (bt *boundaryTag) readBlock(addr byte) block {
	return block{
		addr: addr,
		size: bt.memory[addr],
		flag: bt.memory[addr+1],
	}
}

func (bt *boundaryTag) writeBlock(b block) {
	for i := b.addr; i < b.addr+b.size+metadataSize; i++ {
		bt.memory[i] = 0
	}
	bt.memory[b.addr] = b.size
	bt.memory[b.addr+1] = b.flag
	bt.memory[b.addr+2+b.size] = b.size
}
