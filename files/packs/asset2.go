package packs

import (
	"bytes"
	"compress/zlib"
	"fl-library/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Asset2 struct {
	Path    string
	Name    string
	IsLoose bool

	NameHash   uint64
	Offset     uint64
	RealSize   uint64
	PackedSize uint64
	IsZip      bool
	Crc32      uint32
}

func (a *Asset2) LoadFromBinary(f *os.File) {
	a.IsLoose = false

	utils.ReadUInt64L(f, &a.NameHash)
	utils.ReadUInt64L(f, &a.Offset)
	utils.ReadUInt64L(f, &a.PackedSize)

	var isZip uint32
	utils.ReadUInt32L(f, &isZip)
	if isZip == 0x10 || isZip == 0x00 {
		a.IsZip = false
	} else if isZip == 0x01 || isZip == 0x11 {
		a.IsZip = true
	}
	utils.ReadUInt32L(f, &a.Crc32)
}

func (a *Asset2) UnpackFromBinary(f *os.File, outDir string) {

	// Open asset file
	if a.Name == "" {
		outDir += fmt.Sprintf(string(filepath.Separator)+"%016x.bin", a.NameHash)
	} else {
		outDir += fmt.Sprintf(string(filepath.Separator)+"%s", a.Name)
	}

	file, err := os.OpenFile(outDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	utils.Check(err)
	defer file.Close()

	// Write asset data
	if !a.IsZip {
		buffer := make([]byte, int64(a.PackedSize))
		utils.FileReadAt(f, buffer, int64(a.Offset))

		utils.FileWrite(file, buffer)
	} else {
		utils.FileSeek(f, int64(a.Offset+4), 0) // Skip to zlib data

		var realSize uint32
		utils.ReadUInt32B(f, &realSize)

		var b bytes.Buffer
		r, err := zlib.NewReader(f)
		utils.Check(err)
		defer r.Close()

		_, err = io.Copy(&b, r)
		utils.Check(err)

		utils.FileWrite(file, b.Bytes())
	}
}
