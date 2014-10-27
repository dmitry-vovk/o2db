package db

import (
	. "dbtest"
	"os"
	"testing"
)

func TestIntIndex(t *testing.T) {
	idx := NewIntIndex()
	idx.Add(IntTest1, TestId1)
	idx.Add(IntTest2, TestId2)
	idx.Add(IntTest1, TestId3)
	found1 := idx.Find(IntTest1) // should return []int{37, 132}
	if found1[0] != TestId1 {
		t.Fatal("Finding by string did not work")
	}
	if found1[1] != TestId3 {
		t.Fatal("Finding by string did not work")
	}
	idx.Delete(IntTest1, TestId1)
	found2 := idx.Find(IntTest1) // // should return []int{132}
	if found2[0] != TestId3 {
		t.Fatal("Finding by string did not work")
	}
	// Test file IO
	err := idx.FlushToFile(IndexFile)
	if err != nil {
		t.Fatalf("Error flushing index to file: %s", err)
	}
	defer os.Remove(IndexFile)
	idx2, err := OpenIntIndex(IndexFile)
	if err != nil {
		t.Fatalf("Error reading index from file: %s", err)
	}
	if len(idx.Map) != len(idx2.Map) {
		t.Fatal("Read index not equal to stored")
	}
}
