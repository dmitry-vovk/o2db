package index_string

import (
	. "dbtest"
	"testing"
)

func TestStringIndex(t *testing.T) {
	idx := NewStringIndex(IndexFile)
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
	if found2[0] != TestId1 {
		t.Fatal("Finding by string did not work")
	}
	// found2 should contain two ids with one version each
	for _, k := range found2 {
		if !(k == TestId1 || k == TestId3) {
			t.Fatalf("Finding by string did not work (id) %d", k)
		}
	}
	// Try deleting single value/id/version
	idx.Delete(StringTestValue1, TestId1, 1)
	found3 := idx.Find(StringTestValue1)
	if len(found3) != 1 {
		t.Fatal("Deleting by id did not work")
	}
}
