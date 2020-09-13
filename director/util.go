package director

import (
	"encoding/binary"
	"io"
	"log"
)

func byteOrder(littleEndian bool) binary.ByteOrder {
	if littleEndian {
		return binary.LittleEndian
	} else {
		return binary.BigEndian
	}
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
	return byteOrder(littleEndian).Uint32(b)
}

func readFourCC(r io.Reader, littleEndian bool) string {
	b := make([]byte, 4)
	n, err := r.Read(b)
	if err != nil || n != 4 {
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
