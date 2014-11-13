package db

import (
	"errors"
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
	for field, cond := range q.Conditions {
		if _, ok := c.Indices[field]; !ok {
			logger.ErrorLog.Printf("Field %s does not have index", field)
			return nil, errors.New("Field " + field + " does not have index")
		}
		var ids map[int][]int
		index := c.Indices[field]
		switch cond.(type) {
		case []interface{}:
			logger.ErrorLog.Printf("Field %s cond is LIST", field)
			// TODO
		case string:
			ids = index.Find(cond)
		case float64:
			logger.ErrorLog.Printf("Field %s cond is EXACT", field)
			switch index.(type) {
			case *IntIndex:
				logger.ErrorLog.Printf("Field %s has IntIndex", field)
				ids = index.Find(cond)
			case *FloatIndex:
				logger.ErrorLog.Printf("Field %s has FloatIndex", field)
				ids = index.Find(cond)
			default:
				logger.ErrorLog.Printf("Field %s ", field)
			}
		case map[string]interface{}:
			logger.ErrorLog.Printf("Field %s cond is COMPLEX", field)
			// TODO
		}
		logger.ErrorLog.Printf("Found in %s index: %# v", field, pretty.Formatter(ids))
	}
	return nil, nil
}
