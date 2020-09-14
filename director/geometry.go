package director

import (
	"fmt"
)

type Rect struct {
	// https://docs.google.com/document/d/1jDBXE4Wv1AEga-o1Wi8xtlNZY4K2fHxW2Xs8RgARrqk/edit#heading=h.81zcuo2hadvj
	// Maybe x1, y1, x2, y2

	A int16
	B int16
	C int16
	D int16
}

func (r Rect) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v]", r.A, r.B, r.C, r.D)
}
