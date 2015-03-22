// This package contains stuff used in tests only.
package dbtest

import (
	. "types"
)

const (
	TestDataDir        = "/tmp/dbtest"
	TestCollectionName = "new / collection"
	TestCollectionHash = "e2b34345c73ebfae7ae71c940dca37746a987ca2"
	TestDataFileName   = TestDataDir + "/" + TestCollectionHash + "/" + DataFileName
	TestIndexFileName  = TestDataDir + "/" + TestCollectionHash + "/" + ObjectIndexFileName
)
