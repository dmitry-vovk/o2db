package index_int

import (
	. "dbtest"
	"os"
	"testing"
)

func TestIntIndex(t *testing.T) {
	idx := NewIntIndex(IndexFile)
	idx.Add(IntTestValue1, TestId1, 0)
	idx.Add(IntTestValue2, TestId2, 0)
	idx.Add(IntTestValue1, TestId3, 0)
	// Try to find non-existing version
	found0 := idx.Find(IntTestValue1)
	if len(found0) != 2 {
		t.Fatal("Finding non existing int value did not work")
	}
	// Try to find non-existing value
	found1 := idx.Find(IntTestValue3)
	if len(found1) != 0 {
		t.Fatal("Finding non existing int value did not work")
	}
	// Try to find existing value and version
	found2 := idx.Find(IntTestValue1)
	if found2[0] != TestId1 {
		t.Fatal("Finding by int did not work")
	}
	// found2 should contain two ids with one version each
	for _, k := range found2 {
		if !(k == TestId1 || k == TestId3) {
			t.Fatalf("Finding by int did not work (id) %d", k)
		}
	}
	// Try deleting single value/id/version
	idx.Delete(IntTestValue1, TestId1, 0)
	found3 := idx.Find(IntTestValue1)
	if len(found3) != 1 {
		t.Fatal("Deleting by id did not work")
	}
	// Test file IO
	err := idx.FlushToFile()
	if err != nil {
		t.Fatalf("Error flushing index to file: %s", err)
	}
	defer os.Remove(IndexFile)
}
