package packs

import (
	"fl-library/utils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Pack struct {
	Path string
	Name string

	Assets []Asset
}

func (p *Pack) LoadFromFile(path string) {

	// Grab absolute path and check for extension
	path, err := filepath.Abs(path)
	utils.Check(err)
	if filepath.Ext(path) != p.getExt() {
		log.Println("Mismatch file extension")
		return
	}

	// Open pack file
	file, err := os.Open(path)
	utils.Check(err)
	defer file.Close()

	// Store basic info
	p.Name = strings.TrimSuffix(filepath.Base(path), p.getExt())
	p.Path = path

	// Load pack
	fmt.Printf("Loading '%s'..\n", path)

	var nextChunk, assetAmount uint32
	chunkCount := 0
	for {

		// Read chunk info
		utils.ReadUInt32B(file, &nextChunk)
		utils.ReadUInt32B(file, &assetAmount)

		// Load assets
		fmt.Printf("Loading %d assets from chunk %d..\n", assetAmount, chunkCount)
		var a Asset
		for i := 0; i < int(assetAmount); i++ {
			a.Path = p.Path
			a.LoadFromBinary(file)

			p.Assets = append(p.Assets, a)
		}

		chunkCount++
		file.Seek(int64(nextChunk), 0)

		if nextChunk == 0 {
			println("Finished!\n")
			break
		}
	}
}

func (p *Pack) LoadFromDir(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	// Load files
	fmt.Printf("Loading '%s' as pack..\n", path)

	files, err := ioutil.ReadDir(path)
	utils.Check(err)

	var a Asset
	for _, f := range files {
		a.Path = path + `\` + f.Name()
		a.IsLoose = true
		a.Name = []byte(f.Name())
		a.Size = uint32(f.Size())

		// Generate checksum
		file, err := os.Open(a.Path)
		utils.Check(err)
		defer file.Close()

		buffer := make([]byte, a.Size)
		file.Read(buffer)
		a.Crc32 = utils.CalcCrc32(buffer)

		p.Assets = append(p.Assets, a)
	}

	println("Finished!\n")
}

func (p *Pack) Unpack(outDir string) {

	// Grab absolute path
	outDir, err := filepath.Abs(outDir)
	utils.Check(err)

	// Create output directory
	outDir += `\` + p.Name
	err = os.MkdirAll(outDir, 0666)
	utils.Check(err)

	// Open pack file
	file, err := os.Open(p.Path)
	utils.Check(err)
	defer file.Close()

	// Unpack assets
	fmt.Printf("Unpacking %s..\n", p.Path)

	for _, a := range p.Assets {
		a.UnpackFromBinary(file, outDir)
	}

	println("Finished!\n")
}

func (p *Pack) WritePack(outDir, outName string) {

	// Grab absolute path
	outDir, err := filepath.Abs(outDir)
	utils.Check(err)

	// Create dir
	err = os.MkdirAll(outDir, 0666)
	utils.Check(err)

	// Take this one chunk at a time
	outFile, err := os.OpenFile(outDir+`\`+outName+".pack", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	utils.Check(err)
	defer outFile.Close()

	fileCount, chunkCount := 0, 0
	for {

		// Store chunk offset
		chunkOffset := utils.Tell(outFile)
		dataOffset := chunkOffset + p.getPadding()

		// Write chunk padding
		outFile.Write(make([]byte, p.getPadding()))
		outFile.Seek(chunkOffset, 0)

		// Write chunk info
		utils.WriteUInt32B(outFile, 0) // NextChunk dummy
		utils.WriteUInt32B(outFile, 0) // FileAmount dummy

		// Iterate through assets
		chunkFileAmount := 0
		for _, a := range p.Assets[fileCount:] {

			// Avoid using all padding space
			if utils.Tell(outFile)+a.GetSize() >= chunkOffset+p.getPadding() {
				break
			}

			utils.WriteUInt32B(outFile, uint32(len(a.Name)))
			outFile.Write([]byte(a.Name))
			utils.WriteUInt32B(outFile, uint32(dataOffset))
			utils.WriteUInt32B(outFile, a.Size)
			utils.WriteUInt32B(outFile, a.Crc32)

			// Open asset file
			inFile, err := os.Open(a.Path)
			utils.Check(err)
			defer inFile.Close()

			// Write asset data to pack
			buffer := make([]byte, a.Size)
			inFile.Read(buffer)
			outFile.WriteAt(buffer, dataOffset)
			dataOffset += int64(a.Size)
			chunkFileAmount++

			// Check for leftover files
			if fileCount+chunkFileAmount == len(p.Assets) {
				dataOffset = 0
				break
			}
		}
		fileCount += chunkFileAmount

		// Write next chunk offset
		// TODO: Do that in a loop
		outFile.Seek(chunkOffset, 0)
		utils.WriteUInt32B(outFile, uint32(dataOffset))

		// Write chunk file amount
		utils.WriteUInt32B(outFile, uint32(chunkFileAmount))
		outFile.Seek(dataOffset, 0)

		// Check for remaining files
		chunkCount++
		if fileCount == len(p.Assets) {
			break
		}
	}
}

func (p *Pack) getExt() string {
	return ".pack"
}

func (p *Pack) getPadding() int64 {
	return 0x2000
}

// Display shows pack info
func (p Pack) Display() {
	fmt.Printf("NAME\t'%s'\nPATH\t'%s'\nAMOUNT\t%d\n\n", p.Name, p.Path, len(p.Assets))
}
