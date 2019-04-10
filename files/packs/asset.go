package packs

import (
	"fl-library/utils"
	"fmt"
	"os"
	"strconv"
)

type Asset struct {
	Path    string
	Name    []byte
	IsLoose bool

	Offset uint32
	Size   uint32
	Crc32  uint32
}

func (a *Asset) LoadFromBinary(f *os.File) {
	var nameLength uint32

	a.IsLoose = false
	utils.ReadUInt32B(f, &nameLength)
	a.Name = make([]byte, int(nameLength))
	f.Read(a.Name)

	utils.ReadUInt32B(f, &a.Offset)
	utils.ReadUInt32B(f, &a.Size)
	utils.ReadUInt32B(f, &a.Crc32)
}

func (a *Asset) UnpackFromBinary(f *os.File, outDir string) {

	// Open asset file
	outDir += `\` + string(a.Name)
	file, err := os.OpenFile(outDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	utils.Check(err)
	defer file.Close()

	// Write data to file
	buffer := make([]byte, int64(a.Size))
	f.ReadAt(buffer, int64(a.Offset))
	file.Write(buffer)
}

// GetSize returns size off asset header
func (a *Asset) GetSize() int64 {
	return int64(4 + len(a.Name) + 4 + 4 + 4)
}

// Display shows pack info
func (a Asset) Display() {
	fmt.Printf("PATH\t'%s'\nLOOSE\t%s\n", a.Path, strconv.FormatBool(a.IsLoose))
	fmt.Printf("NAME\t'%s'\nOFFSET\t%#x\nSIZE\t%d\nCRC32\t%#010x\n\n", a.Name, a.Offset, a.Size, a.Crc32)
}
