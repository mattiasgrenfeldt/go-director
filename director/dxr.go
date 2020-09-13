package director

import (
	"fmt"
	"io"
)

type DXR struct{}

func Parse(r io.Reader) {
	rifx := parseRifx(r)

	c := rifx.chunks[0]
	fmt.Printf("%v %v %v\n", c.fourCC, len(c.data), c.data)
}
