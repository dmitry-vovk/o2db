package db

import (
	"os"
	"testing"
)

const (
	stringTest1 = "hello, world"
	stringTest2 = "... and good bye!"
	testId1     = 37
	testId2     = 50
	testId3     = 132
	indexFile   = "/tmp/strings.index.tmp"
)

func TestStringIndex(t *testing.T) {
	idx := NewStringIndex()
	idx.Add(stringTest1, testId1)
	idx.Add(stringTest2, testId2)
	idx.Add(stringTest1, testId3)
	found1 := idx.Find(stringTest1) // should return []int{37, 132}
	if found1[0] != testId1 {
		t.Fatal("Finding by string did not work")
	}
	if found1[1] != testId3 {
		t.Fatal("Finding by string did not work")
	}
	idx.Delete(stringTest1, testId1)
	found2 := idx.Find(stringTest1) // // should return []int{132}
	if found2[0] != testId3 {
		t.Fatal("Finding by string did not work")
	}
	// Test file IO
	err := idx.FlushToFile(indexFile)
	if err != nil {
		t.Fatalf("Error flushing index to file: %s", err)
	}
	defer os.Remove(indexFile)
	idx2, err := OpenStringIndex(indexFile)
	if err != nil {
		t.Fatalf("Error reading index from file: %s", err)
	}
	if len(idx.Map) != len(idx2.Map) {
		t.Fatal("Read index not equal to stored")
	}
}
