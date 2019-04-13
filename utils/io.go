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

func FileSeek(f *os.File, offset int64, whence int) {
	_, err := f.Seek(offset, whence)
	Check(err)
}

// Read
func FileRead(f *os.File, b []byte) {
	_, err := f.Read(b)
	Check(err)
}

func FileReadAt(f *os.File, b []byte, off int64) {
	_, err := f.ReadAt(b, off)
	Check(err)
}

// Write
func FileWrite(f *os.File, b []byte) {
	_, err := f.Write(b)
	Check(err)
}

func FileWriteAt(f *os.File, b []byte, off int64) {
	_, err := f.WriteAt(b, off)
	Check(err)
}

func FileWriteString(f *os.File, s string) {
	_, err := f.WriteString(s)
	Check(err)
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
