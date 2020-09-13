package director

import (
	"io"
	"log"
)

const bigEndianMagic = "RIFX"
const littleEndianMagic = "XFIR"

type rifx struct {
	littleEndian    bool
	directorVersion string
	chunks          []*rifxChunk
}

type rifxChunk struct {
	fourCC string
	data   []byte
}

func parseRifx(r io.Reader) *rifx {
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

	read := uint32(4)
	var chunks []*rifxChunk
	for read != size {
		c := parseRifxChunk(r, le)
		n := len(c.data)
		read += uint32(n + (n % 2) + 8) // (c.size % 2) is for extra pad byte.
		chunks = append(chunks, c)
	}

	return &rifx{littleEndian: le, directorVersion: version, chunks: chunks}
}

func parseRifxChunk(r io.Reader, littleEndian bool) *rifxChunk {
	fourCC := readFourCC(r, littleEndian)
	size := readUint32(r, littleEndian)
	data := make([]byte, size)
	n, err := r.Read(data)
	if err != nil || n != int(size) {
		log.Fatalf("parseRifxChunk bad reader err: %v n: %v\n", err, n)
	}
	if size%2 != 0 {
		// Odd size, read one pad byte.
		n, err := r.Read([]byte{0})
		if err != nil || n != 1 {
			log.Fatalf("parseRifxChunk failed to read pad byte, err: %v n: %v\n", err, n)
		}
	}
	return &rifxChunk{fourCC: fourCC, data: data}
}

func readInt32(r io.Reader, littleEndian bool) int32 {
	return int32(readUint32(r, littleEndian))
}

func readUint32(r io.Reader, littleEndian bool) uint32 {
	b := make([]byte, 4)
	n, err := r.Read(b)
	if err != nil || n != 4 {
		log.Fatalf("readUint32 bad read err: %v n: %v\n", err, n)
	}
	var x uint32
	if littleEndian {
		x = (uint32(b[3]) << 24) + (uint32(b[2]) << 16) + (uint32(b[1]) << 8) + uint32(b[0])
	} else {
		x = (uint32(b[0]) << 24) + (uint32(b[1]) << 16) + (uint32(b[2]) << 8) + uint32(b[3])
	}
	return x
}

func readFourCC(r io.Reader, littleEndian bool) string {
	b := make([]byte, 4)
	n, err := r.Read(b)
	if err != nil || n != 4 {
		//panic("bad")
		log.Fatalf("readFourCC bad read err: %v n: %v\n", err, n)
	}
	if littleEndian {
		for i := 0; i < 2; i++ {
			j := n - i - 1
			b[i], b[j] = b[j], b[i]
		}
	}
	return string(b)
}
