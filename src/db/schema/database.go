package schema

type Database struct {
	DataDir string
	Collections map[string]Collection
}
