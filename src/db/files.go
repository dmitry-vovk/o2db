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
	var err error
	dbFile.Handler, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return nil, err
	}
	// We are not trying to map empty file
	if stat, _ := os.Stat(fileName); stat.Size() == 0 {
		dbFile.Handler.Truncate(int64(os.Getpagesize()))
	}
	dbFile.Buffer, err = mmap.Map(dbFile.Handler, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &dbFile, nil
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
	if len(this.Buffer) < offset+len(data) {
		// TODO increase length/size of slice/file
		return ErrNotImplemented
	}
	if n := copy(this.Buffer[offset:], data); n != len(data) {
		return ErrTruncated
	}
	return nil
}
