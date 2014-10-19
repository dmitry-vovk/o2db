package db

import (
	"bytes"
	"errors"
	"logger"
	"os"
)

var (
	ErrTruncated      = errors.New("Could not write all data")
	ErrNotImplemented = errors.New("TODO!")
)

type DbFile struct {
	FileName string
	Handler  *os.File
}

func (f *DbFile) Open() {
	if err := f.openFile(); err != nil {
		logger.ErrorLog.Printf("Error opening data file: %s", err)
	}
}

// Opens a file and maps it into memory
func NewFile(fileName string) (*DbFile, error) {
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
	return err
}

// Flush, unmap, and close the file
func (f *DbFile) Close() error {
	if f.Handler == nil {
		return errors.New("File is not open")
	}
	return f.Handler.Close()
}

// Return portion of file starting at 'offset' and of 'len' length
func (f *DbFile) Read(offset, count int) ([]byte, error) {
	buf := make([]byte, count)
	_, err := f.Handler.ReadAt(buf, int64(offset))
	return buf, err
}

// Write 'data' bytes starting at 'offset'
func (f *DbFile) Write(data []byte, offset int) error {
	_, err := f.Handler.WriteAt(data, int64(offset))
	return err
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
