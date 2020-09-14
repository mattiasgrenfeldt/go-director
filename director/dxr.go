package director

import (
	"fmt"
	"log"
	"os"
)

type DXR struct{}

func ParseDXR(r *os.File) DXR {
	rifx := ParseRifx(r)

	imap := ParseImap(r, rifx.Chunks[0])
	//fmt.Printf("%v %v %v %v\n", imap.MmapCount, imap.MmapPos, imap.MmapVersion, imap.Unknown)

	c1 := rifx.Chunks[1]
	if imap.MmapPos != c1.Offset {
		log.Fatalf("ParseDXR mmap is not second chunks, pos got: %v want: %v", imap.MmapPos, c1.Offset)
	}
	mmap := ParseMmap(r, c1)
	fmt.Printf("-- Mmap\n%v\n", mmap)

	for i, r := range mmap.Resources {
		fmt.Printf("-- Res %d\n%v\n", i, r)
	}

	return DXR{}
}
