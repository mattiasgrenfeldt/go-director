package director

import (
	"io"
	"log"
)

const bigEndianMagic = "RIFX"
const littleEndianMagic = "XFIR"

// rifx has the following structure:
//  RIFX or XFIR    - string  - 4 bytes - this decides whether the rest of the file is bigEndian or littleEndian respectively. This even controls whether fourCCs are reversed or not.
//  size            - uint32  - 4 bytes - amount of data to follow.
//  directorVersion - string  - 4 bytes - Example "MV93".
//  chunks          - []chunk - (size-4) bytes
type rifx struct {
	littleEndian    bool
	directorVersion string
	chunks          []rifxChunk
}

// rifxChunk has the following structure:
//  fourCC - string - 4 bytes
//  size   - uint32 - 4 bytes
//  data   - []byte - size bytes
type rifxChunk struct {
	littleEndian bool
	fourCC       string
	size         uint32
	offset       int64
}

func parseRifx(r io.ReadSeeker) rifx {
	magic := readFourCC(r, false)
	var le bool
	if magic == bigEndianMagic {
		le = false
	} else if magic == littleEndianMagic {
		le = true
	} else {
		panic("Bad magic")
	}

	size := int64(readUint32(r, le))
	version := readFourCC(r, le)

	offset := int64(12)
	var chunks []rifxChunk
	for offset != size+8 {
		c := parseRifxChunk(r, le, offset)
		offset += int64(c.size + (c.size % 2) + 8) // (c.size % 2) is for extra pad byte.
		chunks = append(chunks, c)
	}
	_, err := r.Read([]byte{1})
	if err != io.EOF {
		log.Fatalf("parseRifx: More data at end of file")
	}

	return rifx{littleEndian: le, directorVersion: version, chunks: chunks}
}

func parseRifxChunk(r io.ReadSeeker, littleEndian bool, offset int64) rifxChunk {
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
	return rifxChunk{littleEndian: littleEndian, fourCC: fourCC, size: size, offset: offset}
}
