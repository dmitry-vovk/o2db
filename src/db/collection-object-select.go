package db

import (
	"github.com/kr/pretty"
	"logger"
	. "types"
)

/*
   Conditions: {
       "prop1": []interface {}{
           float64(1),
           float64(2),
           float64(5),
       },
       "prop2": float64(4),
       "prop3": map[string]interface {}{
           "<":  float64(2.5),
           ">=": float64(1),
       },
       "prop4": "hello"
   },
*/

func (c *Collection) SelectObjects(q SelectObjects) ([]*ObjectFields, error) {
	logger.ErrorLog.Printf("%# v", pretty.Formatter(q))
	var result []*ObjectFields
	//var result []*ObjectFields
	/*
		for field, cond := range q.Query {
			if c.isScalarValue(cond) {
				logger.ErrorLog.Printf("Value %s, %# v", field, pretty.Formatter(cond))
			} else {
				logger.ErrorLog.Printf("Expression %s, %# v", field, pretty.Formatter(cond))
			}
		}*/
	c.processQuery("", q.Query, "")
	/*
		if _, ok := c.Indices[field]; !ok {
			logger.ErrorLog.Printf("Field %s does not have index", field)
			return nil, errors.New("Field " + field + " does not have index")
		}
		var ids map[int][]int
		index := c.Indices[field]
		switch cond.(type) {
		case []interface{}:
			logger.ErrorLog.Printf("Field %s cond is LIST", field)
			for val := range cond.([]interface{}) {
				result = append(result, c.collectFoundObjects(index.Find(val))...)
			}
		case string:
			result = append(result, c.collectFoundObjects(index.Find(cond))...)
		case float64:
			logger.ErrorLog.Printf("Field %s cond is EXACT", field)
			switch index.(type) {
			case *IntIndex:
				logger.ErrorLog.Printf("Field %s has IntIndex", field)
				result = append(result, c.collectFoundObjects(index.Find(cond))...)
			case *FloatIndex:
				logger.ErrorLog.Printf("Field %s has FloatIndex", field)
				result = append(result, c.collectFoundObjects(index.Find(cond))...)
			}
		case map[string]interface{}:
			logger.ErrorLog.Printf("Field %s cond is COMPLEX", field)
			for op, val := range cond.(map[string]interface{}) {
				result = append(result, c.collectFoundObjects(index.ConditionalFind(op, val))...)
			}
		}
		logger.ErrorLog.Printf("Found in %s index: %# v", field, pretty.Formatter(ids))
	}*/
	return result, nil
}

func (c *Collection) processQuery(verb string, q ObjectFields, indent string) map[int][]int {
	var foundIds map[int][]int
	/*
		ids := make(map[int][]int)
		for field, cond := range q {
			if c.isConditional(field) {
				ids[len(ids)] = c.processQuery(field, cond.(map[string]interface{}), indent+"  ")
			} else if c.isScalarValue(cond) {
				logger.ErrorLog.Printf("%sField %s = %v", indent, field, cond)
				ids[len(ids)] = c.Indices[field].Find(cond)
			} else {
				logger.ErrorLog.Printf("%s%s", indent, field)
				ids[len(ids)] = c.processQuery(field, cond.(map[string]interface{}), indent+"  ")
			}
		}
		switch verb {
		case "OR":
		case "AND":
		case "XOR":
		case "NOT":
		default:
			// verb undefined
		}
		logger.ErrorLog.Printf("Found %# v", pretty.Formatter(ids))
	*/
	return foundIds
}

// Return all
func (c *Collection) joinOR(ids []map[int][]int) map[int][]int {
	result := make(map[int][]int)
	return result
}

// Return only those present in all slices
func (c *Collection) joinAND(ids []map[int][]int) map[int][]int {
	result := make(map[int][]int)
	return result
}

// Return only ones that are unique in slices
func (c *Collection) joinXOR(ids []map[int][]int) map[int][]int {
	result := make(map[int][]int)
	return result
}

// Return all those present in first but not is the last
func (c *Collection) joinNOT(ids []map[int][]int) map[int][]int {
	result := make(map[int][]int)
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

// TODO wrong but way to go
func (c *Collection) collectFoundObjects(ids map[int][]int) []*ObjectFields {
	var foundObjects []*ObjectFields
	for id, versions := range ids {
		objects, err := c.getObjectByIdAndVersion(id, versions[len(versions)-1])
		if err == nil {
			foundObjects = append(foundObjects, objects)
		}
	}
	return foundObjects
}
