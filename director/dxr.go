package director

import (
	"fmt"
	"os"
)

type DXR struct{}

func ParseDXR(r *os.File) DXR {
	rifx := parseRifx(r)

	/*
		for _, c := range rifx.chunks {
			fmt.Printf("%q %v %v\n", c.fourCC, c.offset, c.size)
		}
	*/

	imap := ParseImap(r, rifx.chunks[0])
	fmt.Printf("%v %v %v %v\n", imap.MemMapCount, imap.MemMapPos, imap.MemMapVersion, imap.Unknown)

	return DXR{}
}
