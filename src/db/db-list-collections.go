package db

func (d *Database) ListCollections() []string {
	var cList []string
	for cName, _ := range d.Collections {
		cList = append(cList, cName)
	}
	return cList
}
