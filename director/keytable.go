package director

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const keyTableFourCC = "KEY*"
const keyTableHeaderLen = 2*2 + 4*2

type KeyTable struct {
	// https://docs.google.com/document/d/1jDBXE4Wv1AEga-o1Wi8xtlNZY4K2fHxW2Xs8RgARrqk/edit#heading=h.t8soymqf9e9h

	// PropertiesLength seems to always be 12.
	PropertiesLength int16
	// KeyLength seems to always be 12.
	KeyLength       int16
	KeyLengthMax    int32
	KeyLengthUnused int32
	Table           []KeyElement
}

func (t KeyTable) String() string {
	return fmt.Sprintf(`PropertiesLength: %v
KeyLength:        %v
KeyLengthMax:     %v
KeyLengthUnused:  %v
`, t.PropertiesLength, t.KeyLength, t.KeyLengthMax, t.KeyLengthUnused)
}

const keyElementLen = 4 * 3

type KeyElement struct {
	OwnedResID int32
	OwnerResID int32
	// FourCC is the FourCC of the owned resource.
	FourCC string
}

func (e KeyElement) String() string {
	return fmt.Sprintf(`OwnedResID: %v
OwnerResID: %v
FourCC:     %q
`, e.OwnedResID, e.OwnerResID, e.FourCC)
}

func ParseKeyTable(r io.ReadSeeker, c RifxChunk) KeyTable {
	if c.FourCC != keyTableFourCC {
		log.Fatalf("ParseKeyTable fourCC got: %v want: %v", c.FourCC, keyTableFourCC)
	}

	if n := (c.Size - keyTableHeaderLen) % keyElementLen; n != 0 {
		log.Fatalf("ParseKeyTable bad c.Size: %v (c.Size - keyTableHeaderLen) %% keyElementLen: %v", c.Size, n)
	}
	_, err := r.Seek(int64(c.Offset)+8, io.SeekStart)
	if err != nil {
		log.Fatalf("ParseKeyTable got err: %v", err)
	}

	h := struct {
		PropertiesLength int16
		KeyLength        int16
		KeyLengthMax     int32
		KeyLengthUnused  int32
	}{}

	err = binary.Read(r, byteOrder(c.LittleEndian), &h)
	if err != nil {
		log.Fatalf("ParseKeyTable got err: %v", err)
	}

	var table []KeyElement
	n := (c.Size - keyTableHeaderLen) / keyElementLen
	for i := uint32(0); i < n; i++ {
		table = append(table, ParseKeyElement(r, c.LittleEndian))
	}

	return KeyTable{
		PropertiesLength: h.PropertiesLength,
		KeyLength:        h.KeyLength,
		KeyLengthMax:     h.KeyLengthMax,
		KeyLengthUnused:  h.KeyLengthUnused,
		Table:            table,
	}
}

func ParseKeyElement(r io.Reader, littleEndian bool) KeyElement {
	return KeyElement{
		OwnedResID: readInt32(r, littleEndian),
		OwnerResID: readInt32(r, littleEndian),
		FourCC:     readFourCC(r, littleEndian),
	}
}

func (t KeyTable) Lookup(owner int32, fourCC string) int32 {
	owned := int32(-1)
	for _, e := range t.Table {
		if e.OwnerResID == owner && e.FourCC == fourCC {
			if owned != -1 {
				panic("Already owned")
			}
			owned = e.OwnedResID
		}
	}
	return owned
}
