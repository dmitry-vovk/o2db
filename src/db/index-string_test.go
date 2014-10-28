package db

import (
	. "dbtest"
	"os"
	"testing"
)

func TestStringIndex(t *testing.T) {
	idx := NewStringIndex()
	idx.Add(StringTestValue1, TestId1, 1)
	idx.Add(StringTestValue2, TestId2, 1)
	idx.Add(StringTestValue1, TestId3, 1)
	// Try to find non-existing version
	found0 := idx.Find(StringTestValue1)
	if len(found0) != 2 {
		t.Fatal("Finding non existing string value did not work")
	}
	// Try to find non-existing value
	found1 := idx.Find(StringTestValue3)
	if len(found1) != 0 {
		t.Fatal("Finding non existing string value did not work")
	}
	// Try to find existing value and version
	found2 := idx.Find(StringTestValue1)
	if found2[TestId1][0] != 1 {
		t.Fatal("Finding by string did not work")
	}
	// found2 should contain two ids with one version each
	for k, v := range found2 {
		if !(k == TestId1 || k == TestId3) {
			t.Fatal("Finding by string did not work (id) %d", k)
		}
		if len(v) != 1 {
			t.Fatal("Finding by string did not work (version)")
		}
	}
	// Try deleting single value/id/version
	idx.Delete(StringTestValue1, TestId1, 1)
	found3 := idx.Find(StringTestValue1)
	if len(found3) != 1 {
		t.Fatal("Deleting by id did not work")
	}
	// Test file IO
	err := idx.FlushToFile(IndexFile)
	if err != nil {
		t.Fatalf("Error flushing index to file: %s", err)
	}
	defer os.Remove(IndexFile)
	idx2, err := OpenStringIndex(IndexFile)
	if err != nil {
		t.Fatalf("Error reading index from file: %s", err)
	}
	if len(idx.Map) != len(idx2.Map) {
		t.Fatal("Read index not equal to stored")
	}
}
