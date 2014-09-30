package db

import (
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

func (this *DbFile) openFile() error {
	var err error
	this.Handler, err = os.OpenFile(this.FileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}
	// We are not trying to map empty file
	if stat, _ := os.Stat(this.FileName); stat.Size() == 0 {
		this.Handler.Truncate(int64(os.Getpagesize() * 1024 * 32))
	}
	this.Buffer, err = mmap.Map(this.Handler, mmap.RDWR, 0)
	if err != nil {
		return err
	}
	return nil
}

// Flush, unmap, and close the file
func (this *DbFile) Close() error {
	if this.Handler == nil {
		return errors.New("File is not open")
	}
	err := this.Buffer.Flush()
	if err != nil {
		return err
	}
	err = this.Buffer.Unmap()
	if err != nil {
		return err
	}
	err = this.Handler.Close()
	if err != nil {
		return err
	}
	return nil
}

// Return portion of file starting at 'start' and of 'len' length
func (this *DbFile) Read(start, len int) ([]byte, error) {
	return this.Buffer[start : start+len], nil
}

// Write 'data' bytes starting at 'offset'
func (this *DbFile) Write(data []byte, offset int) error {
	var diff = offset + len(data) - len(this.Buffer)
	if diff > 0 {
		var zeroes = make([]byte, diff)
		this.append(zeroes)
	}
	if n := copy(this.Buffer[offset:], data); n != len(data) {
		return ErrTruncated
	}
	return nil
}

// Append data to the end of the file (resizing it)
func (this *DbFile) append(data []byte) error {
	err := this.Close()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(this.FileName, os.O_APPEND|os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return err
	}
	return this.openFile()
}
