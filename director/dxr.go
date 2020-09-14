package director

import (
	"fmt"
	"log"
	"os"
)

const MovieResID = 1024

type DXR struct{}

func ParseDXR(r *os.File) DXR {
	rifx := ParseRifx(r)

	imap := ParseImap(r, rifx.Chunks[0])

	c1 := rifx.Chunks[1]
	if imap.MmapPos != c1.Offset {
		log.Fatalf("ParseDXR mmap is not second chunks, pos got: %v want: %v", imap.MmapPos, c1.Offset)
	}
	mmap := ParseMmap(r, c1)

	ktRes := mmap.Resources[3]
	c2 := rifx.Chunks[2]
	if ktRes.Offset != c2.Offset {
		log.Fatalf("ParseDXR KEY* is not third chunks, pos got: %v want: %v", ktRes.Offset, c2.Offset)
	}
	keyTable := ParseKeyTable(r, c2)

	cfgResID := keyTable.Lookup(MovieResID, configNewFourCC)
	cfgRes := mmap.Resources[cfgResID]
	cfgIndex, cfgChunk := rifx.OffsetToChunk(cfgRes.Offset)
	if cfgIndex == -1 {
		log.Fatalf("ParseDXR couldn't find a Config")
	}
	config := ParseConfig(r, cfgChunk)

	fmt.Printf("-- Config\n%v\n", config)

	return DXR{}
}
