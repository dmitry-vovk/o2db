package index_float

import (
	. "dbtest"
	"os"
	"testing"
)

func TestFloatIndex(t *testing.T) {
	idx := NewFloatIndex(IndexFile)
	idx.Add(FloatTestValue1, TestId1, 0)
	idx.Add(FloatTestValue2, TestId2, 0)
	idx.Add(FloatTestValue1, TestId3, 0)
	// Try to find non-existing version
	found0 := idx.Find(FloatTestValue1)
	if len(found0) != 2 {
		t.Fatal("Finding non existing int value did not work")
	}
	// Try to find non-existing value
	found1 := idx.Find(FloatTestValue3)
	if len(found1) != 0 {
		t.Fatal("Finding non existing int value did not work")
	}
	// Try to find existing value and version
	found2 := idx.Find(FloatTestValue1)
	//log.Printf("!!!!!!!!!!!!!!!!!!!!! %# v", pretty.Formatter(found2))
	if found2[0] != TestId1 {
		t.Fatal("Finding by float did not work")
	}
	// found2 should contain two ids with one version each
	for _, k := range found2 {
		if !(k == TestId1 || k == TestId3) {
			t.Fatalf("Finding by float did not work (id) %d", k)
		}
	}
	// Try deleting single value/id/version
	idx.Delete(FloatTestValue1, TestId1, 0)
	found3 := idx.Find(FloatTestValue1)
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
