package db

import (
	"github.com/kr/pretty"
	"logger"
	. "types"
)

func (c *Collection) SelectObjects(q SelectObjects) ([]*ObjectFields, error) {
	logger.ErrorLog.Printf("%# v", pretty.Formatter(q))
	var result []*ObjectFields
	foundIds := c.processQuery("", q.Query, "")
	logger.ErrorLog.Printf("Result: %# v", pretty.Formatter(foundIds))
	for _, id := range foundIds {
		object, err := c.ReadObject(ReadObject{
			Fields: ObjectFields{
				FIELD_ID: id,
			},
		})
		if err == nil {
			result = append(result, object)
		}
	}
	return result, nil
}

func (c *Collection) processQuery(verb string, q ObjectFields, indent string) []int {
	var foundIds []int
	ids := make(map[int][]int)
	for field, cond := range q {
		if c.isConditional(field) {
			logger.ErrorLog.Printf("%sConditional %s %s %v", indent, verb, field, cond)
			if set := c.Indices[verb].ConditionalFind(field, cond); set != nil {
				ids[len(ids)] = set
			}
		} else if c.isScalarValue(cond) {
			logger.ErrorLog.Printf("%sField %s = %v", indent, field, cond)
			if set := c.Indices[field].Find(cond); set != nil {
				ids[len(ids)] = set
			}
		} else {
			logger.ErrorLog.Printf("%s%s", indent, field)
			if set := c.processQuery(field, cond.(map[string]interface{}), indent+"  "); set != nil {
				ids[len(ids)] = set
			}
		}
	}
	logger.ErrorLog.Printf("Raw result: %# v", pretty.Formatter(ids))
	switch verb {
	case "OR":
		foundIds = c.joinOR(ids)
	case "AND":
		foundIds = c.joinAND(ids)
	case "XOR":
		foundIds = c.joinXOR(ids)
	case "NOT":
		foundIds = c.joinNOT(ids)
	default:
		foundIds = ids[0]
	}
	return foundIds
}

func (c *Collection) slice2map(values []int) map[int]struct{} {
	result := make(map[int]struct{})
	for _, val := range values {
		result[val] = struct{}{}
	}
	return result
}

func (c *Collection) map2slice(values map[int]struct{}) []int {
	result := []int{}
	for id, _ := range values {
		result = append(result, id)
	}
	return result
}

// Return all
func (c *Collection) joinOR(ids map[int][]int) []int {
	result := []int{}
	for _, set := range ids {
		result = append(result, set...)
	}
	return result
}

// Return only those present in all slices
func (c *Collection) joinAND(ids map[int][]int) []int {
	result := []int{}
	return result
}

// Return only ones that are unique in slices
func (c *Collection) joinXOR(ids map[int][]int) []int {
	result := []int{}
	return result
}

// Return all those present in first but not is the last
func (c *Collection) joinNOT(ids map[int][]int) []int {
	result := ids[0]
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
