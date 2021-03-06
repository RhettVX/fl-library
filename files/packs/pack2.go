package packs

import (
	"bytes"
	"fl-library/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type Pack2 struct {
	Path string
	Name string

	NameList []utils.HashName
	Assets   []Asset2
}

func (p *Pack2) LoadFromFile(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	// Open pack2
	inFile, err := os.Open(path)
	utils.Check(err)
	defer inFile.Close()

	// Store basic info
	p.Path = path
	p.Name = strings.TrimSuffix(filepath.Base(p.Path), p.getExt())

	// Load pack2
	id := make([]byte, 4)
	utils.FileRead(inFile, id)
	if !reflect.DeepEqual(id, p.getID()) {
		log.Println("Wrong file magic")
		return
	}

	// Load header
	var assetAmount uint32
	var fileSize, mapOffset uint64
	utils.ReadUInt32L(inFile, &assetAmount)
	utils.ReadUInt64L(inFile, &fileSize)
	utils.ReadUInt64L(inFile, &mapOffset)

	// Load map
	utils.FileSeek(inFile, int64(mapOffset), 0)
	for i := 0; i < int(assetAmount); i++ {
		var a Asset2
		a.Path = p.Path
		a.LoadFromBinary(inFile)

		p.Assets = append(p.Assets, a)

		if a.NameHash == p.getNameHash() {
			p.NameList = a.ReadNameList(inFile)
		}
		// fmt.Printf("%x : %x : %d : %t\n", a.NameHash, a.Offset, a.PackedSize, a.IsZip)
	}
}

func (p *Pack2) LoadFromDir(path string) {
	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	// Load files
	fmt.Printf("Loading '%s' as pack2...\n", path)

	namePath := filepath.Join(path, "{NAMELIST}")
	if _, err := os.Stat(namePath); err == nil {
		err = os.Remove(namePath)
		utils.Check(err)
	} else if os.IsNotExist(err) {

	} else {
		utils.Check(err)
	}

	files, err := ioutil.ReadDir(path)
	utils.Check(err)

	// Generate NameList
	nameFile, err := os.OpenFile(namePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	utils.Check(err)
	defer nameFile.Close()

	utils.FileWrite(nameFile, []byte("{NAMELIST}\x0a"))
	for _, f := range files {
		utils.FileWrite(nameFile, []byte(f.Name()+"\x0a"))
	}

	var a Asset2
	for _, f := range files {
		a.Path = filepath.Join(path, f.Name())
		a.IsLoose = true
		a.Name = f.Name()
		a.NameHash = utils.Pack2Hash(bytes.ToUpper([]byte(a.Name)))
		a.RealSize = uint64(f.Size())

		// Generate checksum
		file, err := os.Open(a.Path)
		utils.Check(err)
		defer file.Close()

		buffer := make([]byte, a.RealSize)
		utils.FileRead(file, buffer)
		a.Crc32 = utils.CalcCrc32(buffer)

		p.Assets = append(p.Assets, a)
	}
}

func (p *Pack2) Unpack(outDir string) {

	// Create output dir
	outDir += string(filepath.Separator) + p.Name
	err := os.MkdirAll(outDir, 0755)
	utils.Check(err)

	// Open pack2 file
	inFile, err := os.Open(p.Path)
	utils.Check(err)
	defer inFile.Close()

	// Unpack assets
	fmt.Printf("Unpacking %s..\n", p.Path)

	for _, a := range p.Assets {
		if a.Name != "{NAMELIST}" {
			a.UnpackFromBinary(inFile, outDir)
		}
	}
	println("Finished!\n")
}

// WritePack2 is quick and dirty
func (p *Pack2) WritePack2(outDir, outName string) {

	// Grab absolute path
	outDir, err := filepath.Abs(outDir)
	utils.Check(err)

	// Create dir
	err = os.MkdirAll(outDir, 0755)
	utils.Check(err)

	file, err := os.OpenFile(filepath.Join(outDir, outName+p.getExt()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	utils.Check(err)
	defer file.Close()

	utils.FileWrite(file, p.getID())
	utils.WriteUInt32L(file, uint32(len(p.Assets)))
	utils.WriteUInt64L(file, 0)                      // This gets replaced later
	utils.WriteUInt64L(file, p.getTotalSize()+0x200) // TODO: Just go unzipped for now
	utils.WriteUInt64L(file, 256)

	pad := make([]byte, p.getTotalSize()+0x200-0x20)
	utils.FileWrite(file, pad)

	var dataOffset uint64 = 0x200
	for _, a := range p.Assets {
		utils.WriteUInt64L(file, a.NameHash)
		utils.WriteUInt64L(file, dataOffset)
		utils.WriteUInt64L(file, a.RealSize)
		utils.WriteUInt32L(file, 0x10)
		utils.WriteUInt32L(file, a.Crc32)

		// Write data
		inFile, err := os.Open(a.Path)
		utils.Check(err)
		defer inFile.Close()

		buffer := make([]byte, a.RealSize)

		utils.FileRead(inFile, buffer)
		utils.FileWriteAt(file, buffer, int64(dataOffset))

		dataOffset += a.RealSize
	}
	packSize := utils.Tell(file)
	utils.FileSeek(file, 0x8, 0)
	utils.WriteUInt64L(file, uint64(packSize))
}

func (p *Pack2) ApplyHash() {
	for i := range p.Assets {
		p.Assets[i].ApplyName(p.NameList)
	}
}

func (p *Pack2) getTotalSize() (output uint64) {
	for _, a := range p.Assets {
		output += a.RealSize
	}
	return output
}

func (p *Pack2) SortAssets() {
	sort.Slice(p.Assets[:], func(i, j int) bool {
		return p.Assets[i].NameHash < p.Assets[j].NameHash
	})
}

func (p *Pack2) getNameHash() uint64 {
	return 0x4137cc65bd97fd30
}

func (p *Pack2) getExt() string {
	return ".pack2"
}

func (p *Pack2) getID() []byte {
	return []byte{'P', 'A', 'K', 0x1}
}
