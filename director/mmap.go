package director

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const mmapFourCC = "mmap"
const mmapHeaderLen = 2*2 + 4*5 // Everything in Mmap except for []Resources.

type Mmap struct {
	// https://docs.google.com/document/d/1jDBXE4Wv1AEga-o1Wi8xtlNZY4K2fHxW2Xs8RgARrqk/edit#heading=h.sqb7ojd28ngq

	// PropertiesLength seems to always be 24.
	PropertiesLength int16
	// ResourceLength seems to always be 20.
	ResourceLength int16
	// ResourcesLengthMax will be the same as the length of Resources.
	ResourcesLengthMax  int32
	ResourcesLengthUsed int32
	LastJunkResID       int32
	PrevMemMapResID     int32
	LastFreeResID       int32
	Resources           []Resource
}

func (m Mmap) String() string {
	return fmt.Sprintf(`PropertiesLength:    %v
ResourceLength:      %v
ResourcesLengthMax:  %v
ResourcesLengthUsed: %v
LastJunkResID:       %v
PrevMemMapResID:     %v
LastFreeResID:       %v
Number of Resources: %v
`, m.PropertiesLength, m.ResourceLength, m.ResourcesLengthMax, m.ResourcesLengthUsed, m.LastJunkResID, m.PrevMemMapResID, m.LastFreeResID, len(m.Resources))
}

const resourceLen = 4 + 4*3 + 4

type Resource struct {
	FourCC    string
	Size      uint32
	Offset    uint32
	Flags     uint32
	LastResID int32
}

func (r Resource) String() string {
	return fmt.Sprintf(`FourCC:   %q
Size:      %v 0x%x
Offset:    %v 0x%x
Flags:     %v 0x%x
LastResID: %v 0x%x
`, r.FourCC, r.Size, r.Size, r.Offset, r.Offset, r.Flags, r.Flags, r.LastResID, r.LastResID)
}

func ParseMmap(r io.ReadSeeker, c RifxChunk) Mmap {
	if c.FourCC != mmapFourCC {
		log.Fatalf("ParseMmap fourCC got: %v want: %v", c.FourCC, mmapFourCC)
	}

	if n := (c.Size - mmapHeaderLen) % resourceLen; n != 0 {
		log.Fatalf("ParseMmap bad c.Size: %v (c.Size - mmapHeaderLen) %% resourceLen: %v", c.Size, n)
	}
	_, err := r.Seek(int64(c.Offset)+8, io.SeekStart)
	if err != nil {
		log.Fatalf("ParseMmap got err: %v", err)
	}

	h := struct {
		PropertiesLength    int16
		ResourceLength      int16
		ResourcesLengthMax  int32
		ResourcesLengthUsed int32
		LastJunkResID       int32
		PrevMemMapResID     int32
		LastFreeResID       int32
	}{}

	err = binary.Read(r, byteOrder(c.LittleEndian), &h)
	if err != nil {
		log.Fatalf("ParseMmap got err: %v", err)
	}

	var res []Resource
	n := (c.Size - mmapHeaderLen) / resourceLen
	for i := uint32(0); i < n; i++ {
		res = append(res, ParseResource(r, c.LittleEndian))
	}

	return Mmap{
		PropertiesLength:    h.PropertiesLength,
		ResourceLength:      h.ResourceLength,
		ResourcesLengthMax:  h.ResourcesLengthMax,
		ResourcesLengthUsed: h.ResourcesLengthUsed,
		LastJunkResID:       h.LastJunkResID,
		PrevMemMapResID:     h.PrevMemMapResID,
		LastFreeResID:       h.LastFreeResID,
		Resources:           res,
	}
}

func ParseResource(r io.ReadSeeker, littleEndian bool) Resource {
	fourCC := readFourCC(r, littleEndian)
	res := struct {
		Size      uint32
		Offset    uint32
		Flags     uint32
		LastResID int32
	}{}
	err := binary.Read(r, byteOrder(littleEndian), &res)
	if err != nil {
		log.Fatalf("ParseMmap couldn't read next Resource: %v", err)
	}

	return Resource{
		FourCC:    fourCC,
		Size:      res.Size,
		Offset:    res.Offset,
		Flags:     res.Flags,
		LastResID: res.LastResID,
	}
}
