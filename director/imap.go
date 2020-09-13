package director

import (
	"encoding/binary"
	"io"
	"log"
)

const imapFourCC = "imap"

type Imap struct {
	// https://docs.google.com/document/d/1jDBXE4Wv1AEga-o1Wi8xtlNZY4K2fHxW2Xs8RgARrqk/edit#heading=h.dq1rrg8abhxt
	MemMapCount   uint32
	MemMapPos     uint32
	MemMapVersion uint32
	Unknown       [12]byte
}

func ParseImap(r io.ReadSeeker, c rifxChunk) Imap {
	if c.fourCC != imapFourCC {
		log.Fatalf("ParseImap fourCC got: %v want: %v", c.fourCC, imapFourCC)
	}
	var imap Imap
	if s := binary.Size(imap); c.size != uint32(s) {
		log.Fatalf("ParseImap size got: %v want: %v", c.size, s)
	}
	_, err := r.Seek(c.offset+8, io.SeekStart)
	if err != nil {
		log.Fatalf("ParseImap got err: %v", err)
	}
	err = binary.Read(r, byteOrder(c.littleEndian), &imap)
	if err != nil {
		log.Fatalf("ParseImap got err: %v", err)
	}
	return imap
}
