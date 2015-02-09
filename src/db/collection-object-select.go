package db

import (
	"errors"
	mapset "github.com/deckarep/golang-set"
	. "types"
)

func (c *Collection) SelectObjects(q SelectObjects) ([]*ObjectFields, uint, error) {
	var result []*ObjectFields
	foundIds := c.processQuery("", q.Query)
	for _, id := range foundIds {
		object, code, err := c.ReadObject(ReadObject{
			Fields: ObjectFields{
				FIELD_ID: id,
			},
		})
		if err == nil {
			result = append(result, object)
		} else {
			return nil, code, errors.New("Error reading object")
		}
	}
	return result, RNoError, nil
}

func (c *Collection) processQuery(verb string, q ObjectFields) []int {
	if len(q) == 0 {
		return []int{}
	}
	var foundIds []int
	ids := make(map[int][]int)
	for field, cond := range q {
		if c.isConditional(field) {
			if set := c.Indices[verb].ConditionalFind(field, cond); set != nil {
				ids[len(ids)] = set
			}
		} else if c.isScalarValue(cond) {
			if set := c.Indices[field].Find(cond); set != nil {
				ids[len(ids)] = set
			}
		} else {
			if set := c.processQuery(field, cond.(map[string]interface{})); set != nil {
				ids[len(ids)] = set
			}
		}
	}
	switch verb {
	case "OR":
		foundIds = c.joinOR(ids)
	case "AND":
		foundIds = c.joinAND(ids)
	case "XOR":
		foundIds = c.joinXOR(ids)
	case "NOT":
		foundIds = c.joinNOT(ids)
	}
	return foundIds
}

// Return all
func (c *Collection) joinOR(ids map[int][]int) []int {
	andSet := mapset.NewSet()
	for _, set := range ids {
		for _, id := range set {
			andSet.Add(id)
		}
	}
	result := []int{}
	for id := range andSet.Iter() {
		result = append(result, id.(int))
	}
	return result
}

// Return only those present in all slices
func (c *Collection) joinAND(ids map[int][]int) []int {
	sets := make([]mapset.Set, len(ids))
	for s, currentSet := range ids {
		sets[s] = mapset.NewSet()
		for _, id := range currentSet {
			sets[s].Add(id)
		}
	}
	set := mapset.NewSet().Union(sets[0])
	for i := 1; i < len(sets); i++ {
		set = set.Intersect(sets[i])
	}
	result := []int{}
	for id := range set.Iter() {
		result = append(result, id.(int))
	}
	return result
}

// Return only ones that are unique in slices
func (c *Collection) joinXOR(ids map[int][]int) []int {
	sets := make([]mapset.Set, len(ids))
	for s, currentSet := range ids {
		sets[s] = mapset.NewSet()
		for _, id := range currentSet {
			sets[s].Add(id)
		}
	}
	set := mapset.NewSet().Union(sets[0])
	for i := 1; i < len(sets); i++ {
		set = set.Difference(sets[i])
	}
	result := []int{}
	for id := range set.Iter() {
		result = append(result, id.(int))
	}
	return result
}

// Return all those present in first but not is the last
func (c *Collection) joinNOT(ids map[int][]int) []int {
	set := mapset.NewSet()
	for _, id := range ids[0] {
		set.Add(id)
	}
	for i := 1; i < len(ids); i++ {
		for _, id := range ids[i] {
			set.Remove(id)
		}
	}
	result := []int{}
	for id := range set.Iter() {
		result = append(result, id.(int))
	}
	return result
}

func (c *Collection) isConditional(val interface{}) bool {
	if _, ok := val.(string); !ok {
		return false
	}
	switch val.(string) {
	case "<", ">", "<=", ">=":
		return true
	}
	return false
}

func (c *Collection) isScalarValue(val interface{}) bool {
	if _, ok := val.(int); ok {
		return true
	}
	if _, ok := val.(string); ok {
		return true
	}
	if _, ok := val.(float64); ok {
		return true
	}
	return false
}
