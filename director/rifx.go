package director

import (
	"io"
	"log"
)

const bigEndianMagic = "RIFX"
const littleEndianMagic = "XFIR"

// Rifx has the following structure:
//  RIFX or XFIR    - string  - 4 bytes - this decides whether the rest of the file is bigEndian or littleEndian respectively. This even controls whether fourCCs are reversed or not.
//  size            - uint32  - 4 bytes - amount of data to follow.
//  directorVersion - string  - 4 bytes - Example "MV93".
//  chunks          - []chunk - (size-4) bytes
type Rifx struct {
	LittleEndian    bool
	DirectorVersion string
	Chunks          []RifxChunk
}

// RifxChunk has the following structure:
//  fourCC - string - 4 bytes
//  size   - uint32 - 4 bytes
//  data   - []byte - size bytes
type RifxChunk struct {
	LittleEndian bool
	FourCC       string
	Size         uint32
	Offset       uint32
}

func ParseRifx(r io.ReadSeeker) Rifx {
	magic := readFourCC(r, false)
	var le bool
	if magic == bigEndianMagic {
		le = false
	} else if magic == littleEndianMagic {
		le = true
	} else {
		panic("Bad magic")
	}

	size := readUint32(r, le)
	version := readFourCC(r, le)

	offset := uint32(12)
	var chunks []RifxChunk
	for offset != size+8 {
		c := ParseRifxChunk(r, le, offset)
		offset += c.Size + (c.Size % 2) + 8 // (c.size % 2) is for extra pad byte.
		chunks = append(chunks, c)
	}
	_, err := r.Read([]byte{1})
	if err != io.EOF {
		log.Fatalf("parseRifx: More data at end of file")
	}

	return Rifx{LittleEndian: le, DirectorVersion: version, Chunks: chunks}
}

func ParseRifxChunk(r io.ReadSeeker, littleEndian bool, offset uint32) RifxChunk {
	fourCC := readFourCC(r, littleEndian)
	size := readUint32(r, littleEndian)
	_, err := r.Seek(int64(size), io.SeekCurrent)
	if err != nil {
		log.Fatalf("parseRifxChunk got err while seeking: %v", err)
	}
	if size%2 != 0 {
		// Odd size, read one pad byte.
		b := make([]byte, 1)
		n, err := r.Read(b)
		if !(err == io.EOF || (err == nil && n == 1 && b[0] == 0)) {
			log.Fatalf("parseRifxChunk failed to read pad byte, err: %v n: %v b[0]: %v\n", err, n, b[0])
		}
	}
	return RifxChunk{LittleEndian: littleEndian, FourCC: fourCC, Size: size, Offset: offset}
}

func (r Rifx) OffsetToChunk(offset uint32) (index int32, c RifxChunk) {
	for i, c := range r.Chunks {
		if c.Offset == offset {
			return int32(i), c
		}
	}
	return -1, RifxChunk{}
}
