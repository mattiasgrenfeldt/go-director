package director

import (
	"encoding/binary"
	"io"
	"log"
)

const imapFourCC = "imap"

type Imap struct {
	// https://docs.google.com/document/d/1jDBXE4Wv1AEga-o1Wi8xtlNZY4K2fHxW2Xs8RgARrqk/edit#heading=h.dq1rrg8abhxt
	MmapCount uint32
	MmapPos   uint32
	// MmapVersion seems to always be 1223.
	MmapVersion uint32
	Unknown     [12]byte
}

func ParseImap(r io.ReadSeeker, c RifxChunk) Imap {
	if c.FourCC != imapFourCC {
		log.Fatalf("ParseImap fourCC got: %v want: %v", c.FourCC, imapFourCC)
	}
	var imap Imap
	if s := binary.Size(imap); c.Size != uint32(s) {
		log.Fatalf("ParseImap size got: %v want: %v", c.Size, s)
	}
	_, err := r.Seek(int64(c.Offset)+8, io.SeekStart)
	if err != nil {
		log.Fatalf("ParseImap got err: %v", err)
	}
	err = binary.Read(r, byteOrder(c.LittleEndian), &imap)
	if err != nil {
		log.Fatalf("ParseImap got err: %v", err)
	}
	return imap
}
