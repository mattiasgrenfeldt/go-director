package director

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const configOldFourCC = "VWCF" // WV = Video Works
const configNewFourCC = "DRCF" // DR = Director

type Config struct {
	// https://docs.google.com/document/d/1jDBXE4Wv1AEga-o1Wi8xtlNZY4K2fHxW2Xs8RgARrqk/edit#heading=h.hbv9bsqk7jgq

	// Length is always encoded as Big Endian.
	Length int16
	// FileVersion is always encoded as Big Endian.
	// 0x163C if protected (DXR).
	FileVersion int16
	SourceRect  Rect
	// MinMember is alwasy encoded as Big Endian.
	// Obsolete: This info is stored in MCsL
	MinMember int16
	// MaxMember is alwasy encoded as Big Endian.
	// Obsolete: This info is stored in MCsL
	MaxMember int16
	Tempo     byte
	_         byte
	BGColor1  byte
	BGColor2  byte
	_         [6]byte
	BGColor0  byte
	_         [25]byte
	Trial     byte
	_         [15]byte
	// OldDefaultPalette is always encoded as Big Endian.
	//  if (FileVersion <= 0x45D) {
	//		defaultPalette = oldDefaultPalette;
	//	}
	OldDefaultPalette int16
	_                 [6]byte
	// DefaultPalette is always encoded as Big Endian
	DefaultPalette uint32
	_              [4]byte
	// 4 bytes smaller than specified in Google Doc
}

func (c Config) String() string {
	return fmt.Sprintf(`Length: %v
FileVersion:       %v 0x%x
SourceRect:        %v
MinMemeber:        %v
MaxMember:         %v
Tempo:             %v
BGColor:           %v %v %v
Trial:             %v
OldDefaultPalette: %v 0x%x
DefaultPalette:    %v 0x%x
`, c.Length, c.FileVersion, c.FileVersion, c.SourceRect, c.MinMember, c.MaxMember, c.Tempo, c.BGColor0, c.BGColor1, c.BGColor2, c.Trial, c.OldDefaultPalette, c.OldDefaultPalette, c.DefaultPalette, c.DefaultPalette)
}

func ParseConfig(r io.ReadSeeker, c RifxChunk) Config {
	if c.FourCC != configOldFourCC && c.FourCC != configNewFourCC {
		log.Fatalf("ParseConfig fourCC got: %v want: %v or %v", c.FourCC, configOldFourCC, configNewFourCC)
	}
	var cfg Config
	if s := binary.Size(cfg); c.Size != uint32(s) {
		log.Fatalf("ParseConfig size got: %v want: %v", c.Size, s)
	}
	_, err := r.Seek(int64(c.Offset)+8, io.SeekStart)
	if err != nil {
		log.Fatalf("ParseConfig got err: %v", err)
	}
	// Note that all fields in Config are Big Endian. Irrespective of whether the
	// file is a RIFX or an XFIR.
	err = binary.Read(r, binary.BigEndian, &cfg)
	if err != nil {
		log.Fatalf("ParseConfig got err: %v", err)
	}
	return cfg
}
