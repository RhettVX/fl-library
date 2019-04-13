package packs

import (
	"bytes"
	"compress/zlib"
	"fl-library/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
		outDir += fmt.Sprintf(string(filepath.Separator)+"%x.bin", a.NameHash)
	} else {
		outDir += fmt.Sprintf(string(filepath.Separator)+"%s", a.Name)
	}

	file, err := os.OpenFile(outDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	utils.Check(err)
	defer file.Close()

	// Write asset data
	if !a.IsZip {
		buffer := make([]byte, int64(a.PackedSize))
		f.ReadAt(buffer, int64(a.Offset))
		file.Write(buffer)
	} else {
		f.Seek(int64(a.Offset+4), 0) // Skip to zlib data

		var realSize uint32
		utils.ReadUInt32B(f, &realSize)

		var b bytes.Buffer
		r, err := zlib.NewReader(f)
		utils.Check(err)
		defer r.Close()

		io.Copy(&b, r)
		file.Write(b.Bytes())
	}
}

func (a *Asset2) ReadNameList(f *os.File) (nameList []utils.HashName) {
	pos, _ := f.Seek(0, 1)

	if a.IsZip {
		f.Seek(int64(a.Offset+8), 0)
		var b bytes.Buffer
		r, err := zlib.NewReader(f)
		utils.Check(err)
		defer r.Close()

		io.Copy(&b, r)

		names := strings.Split(string(b.Bytes()), "\x0a")
		for _, n := range names {
			if n == "" {
				continue
			}

			upperString := bytes.ToUpper([]byte(n))
			hashCaps := utils.Pack2Hash(upperString)
			nameList = append(nameList, utils.HashName{Hash: hashCaps, Name: n})
		}
	} else { // TODO: I am too tired to remember what should go here
		buffer := make([]byte, a.PackedSize)
		f.ReadAt(buffer, int64(a.Offset))

		names := strings.Split(string(buffer), "\x0a")
		for _, n := range names {
			if n == "" {
				continue
			}

			upperString := bytes.ToUpper([]byte(n))
			hashCaps := utils.Pack2Hash(upperString)
			nameList = append(nameList, utils.HashName{Hash: hashCaps, Name: n})
		}
	}

	f.Seek(pos, 0)
	return nameList
}

func (a *Asset2) ApplyName(nameList []utils.HashName) {
	for _, x := range nameList {

		if a.NameHash == x.Hash {
			a.Name = x.Name
			return
		}
	}
}
