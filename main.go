package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

const memorySize byte = 32

type allocator interface {
	name() string
	init(size byte)
	malloc(size byte)
	free(addr byte)
	printSnapshot()
}

func main() {
	//	run[linear](
	//		call{allocator.malloc, 7},
	//		call{allocator.malloc, 2},
	//		call{allocator.malloc, 3},
	//		call{allocator.malloc, 5},
	//	)

	//	run[freeList](
	//		call{allocator.malloc, 7},
	//		call{allocator.malloc, 2},
	//		call{allocator.malloc, 3},
	//		call{allocator.malloc, 5},
	//		call{allocator.free, 0x07},
	//		call{allocator.free, 0x0c},
	//		call{allocator.free, 0x09},
	//		call{allocator.free, 0x00},
	//	)

	run[boundaryTag](
		call{allocator.malloc, 7},
		call{allocator.malloc, 2},
		call{allocator.malloc, 3},
		call{allocator.malloc, 5},
		call{allocator.free, 0x0c},
		call{allocator.free, 0x17},
		call{allocator.free, 0x11},
		call{allocator.free, 0x02},
	)
}

func run[T linear | freeList | boundaryTag](calls ...call) {
	alloc := any(new(T)).(allocator)
	alloc.init(memorySize)

	fmt.Printf("Memory size: %d\n", memorySize)
	fmt.Printf("Allocator type: %s\n", alloc.name())
	fmt.Printf("Press ENTER to print the next snapshot\n\n")

	fmt.Printf("0. init\n")
	alloc.printSnapshot()

	for i, call := range calls {
		fmt.Scanln()
		fmt.Printf("%d. %s\n", i+1, call)
		call.fn(alloc, call.arg)
		alloc.printSnapshot()
	}
}

type call struct {
	fn  func(allocator, byte)
	arg byte
}

func (c call) String() string {
	switch name := funcName(c.fn); name {
	case "malloc":
		return fmt.Sprintf("%s(%d)", name, c.arg)
	case "free":
		return fmt.Sprintf("%s(%#.2x)", name, c.arg)
	default:
		panic("unreachable")
	}
}

func funcName(fn any) string {
	pc := reflect.ValueOf(fn).Pointer()
	name := runtime.FuncForPC(pc).Name()
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
