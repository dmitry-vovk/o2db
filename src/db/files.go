package db

import (
	mmap "github.com/edsrzf/mmap-go"
	"os"
	"errors"
)

type DbFile struct {
	FileName string
	Handler  *os.File
	Buffer   mmap.MMap
}

// Opens a file and maps it into memory
func OpenFile(fileName string) (*DbFile, error) {
	dbFile := &DbFile{
		FileName: fileName,
	}
	var err error
	dbFile.Handler, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return nil, err
	}
	dbFile.Buffer, err = mmap.Map(dbFile.Handler, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}
	return dbFile, nil
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
