package db

import (
	"bytes"
	"errors"
	mmap "github.com/edsrzf/mmap-go"
	"os"
)

var (
	ErrTruncated      = errors.New("Could not write all data")
	ErrNotImplemented = errors.New("TODO!")
)

type DbFile struct {
	FileName string
	Handler  *os.File
	Buffer   mmap.MMap
}

// Opens a file and maps it into memory
func OpenFile(fileName string) (*DbFile, error) {
	dbFile := DbFile{
		FileName: fileName,
	}
	err := dbFile.openFile()
	return &dbFile, err
}

func (f *DbFile) Touch() error {
	return f.openFile()
}

func (f *DbFile) openFile() error {
	var err error
	f.Handler, err = os.OpenFile(f.FileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	// We are not trying to map empty file
	if stat, _ := os.Stat(f.FileName); stat.Size() == 0 {
		f.Handler.Truncate(int64(os.Getpagesize() * 1024 * 32))
	}
	f.Buffer, err = mmap.Map(f.Handler, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	return nil
}

// Flush, unmap, and close the file
func (f *DbFile) Close() error {
	if f.Handler == nil {
		return errors.New("File is not open")
	}
	err := f.Buffer.Flush()
	if err != nil {
		return err
	}
	err = f.Buffer.Unmap()
	if err != nil {
		return err
	}
	err = f.Handler.Close()
	if err != nil {
		return err
	}
	return nil
}

// Return portion of file starting at 'start' and of 'len' length
func (f *DbFile) Read(start, len int) ([]byte, error) {
	return f.Buffer[start : start+len], nil
}

// Write 'data' bytes starting at 'offset'
func (f *DbFile) Write(data []byte, offset int) error {
	var diff = offset + len(data) - len(f.Buffer)
	if diff > 0 {
		var zeroes = make([]byte, diff)
		f.append(zeroes)
	}
	if n := copy(f.Buffer[offset:], data); n != len(data) {
		return ErrTruncated
	}
	return nil
}

// Append data to the end of the file (resizing it)
func (f *DbFile) append(data []byte) error {
	err := f.Close()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(f.FileName, os.O_APPEND|os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(data); err != nil {
		return err
	}
	return f.openFile()
}

// Replace content of the file with data using temporary file
func (f *DbFile) Dump(data *bytes.Buffer) error {
	if f.Handler != nil {
		err := f.Close()
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(f.FileName+".tmp", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(data.Bytes()); err != nil {
		return err
	}
	err = os.Remove(f.FileName)
	if err != nil {
		return err
	}
	err = os.Rename(f.FileName+".tmp", f.FileName)
	if err != nil {
		return err
	}
	return f.openFile()
}
