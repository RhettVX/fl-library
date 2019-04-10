package utils

import (
	"encoding/binary"
	"io"
	"os"
)

func Tell(f *os.File) int64 {
	pos, err := f.Seek(0, 1)
	Check(err)

	return pos
}

// 32 Big
func ReadUInt32B(r io.Reader, i *uint32) {
	err := binary.Read(r, binary.BigEndian, i)
	Check(err)
}

func WriteUInt32B(w io.Writer, i uint32) {
	err := binary.Write(w, binary.BigEndian, i)
	Check(err)
}

// 32 Little
func ReadUInt32L(r io.Reader, i *uint32) {
	err := binary.Read(r, binary.LittleEndian, i)
	Check(err)
}

func WriteUInt32L(w io.Writer, i uint32) {
	err := binary.Write(w, binary.LittleEndian, i)
	Check(err)
}

// 64 Big
func ReadUInt64B(r io.Reader, i *uint64) {
	err := binary.Read(r, binary.BigEndian, i)
	Check(err)
}

func WriteUInt64B(w io.Writer, i uint64) {
	err := binary.Write(w, binary.BigEndian, i)
	Check(err)
}

// 64 Little
func ReadUInt64L(r io.Reader, i *uint64) {
	err := binary.Read(r, binary.LittleEndian, i)
	Check(err)
}

func WriteUInt64L(w io.Writer, i uint64) {
	err := binary.Write(w, binary.LittleEndian, i)
	Check(err)
}
