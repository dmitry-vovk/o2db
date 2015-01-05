package db

import (
	"github.com/kr/pretty"
	"logger"
	. "types"
)

type Subscription struct {
	Key     string
	Query   ObjectFields
	Clients []*Client
}

// Match the object against subscription's Query
func (s *Subscription) Match(object ObjectFields) bool {
	// No conditions, pass all the object updates
	if len(s.Query) == 0 {
		return true
	}
	logger.ErrorLog.Print("Matching: ")
	logger.ErrorLog.Printf("Object: %# v", pretty.Formatter(object))
	logger.ErrorLog.Printf("Query: %# v", pretty.Formatter(s.Query))
	return s.match("", object, "")
}

func (s *Subscription) match(verb string, object ObjectFields, indent string) bool {
	// go over all queries
	for field, cond := range s.Query {
		if s.isAggregate(field) {
			logger.ErrorLog.Printf("%sAggregate %s | %v", indent, field, cond)
			if outcome := s.applyAggregate(field, object, cond.(map[string]interface{})); !outcome {
				return false
			}
		} else if s.isConditional(field) {
			logger.ErrorLog.Printf("%sConditional %s | %s | %v", indent, verb, field, cond)
			if match := s.applyCondition(field, object, cond.(map[string]interface{})); !match {
				return false
			}
		} else if s.isScalarValue(cond) {
			logger.ErrorLog.Printf("%sField %s = %v", indent, field, cond)
			if object[field] != cond {
				return false
			}
		} else {
			logger.ErrorLog.Printf("%s%s", indent, field)
			if match := s.match(field, cond.(map[string]interface{}), indent+"  "); !match {
				return false
			}
		}
	}
	return true
}

func (s *Subscription) isAggregate(fn string) bool {
	switch fn {
	case "OR", "AND", "NOT", "XOR":
		return true
	}
	return false
}

func (s *Subscription) applyAggregate(fn string, object ObjectFields, cond ObjectFields) bool {
	// TODO
	return true
}

func (s *Subscription) isConditional(val string) bool {
	switch val {
	case "<", ">", "<=", ">=":
		return true
	}
	return false
}

func (s *Subscription) applyCondition(op string, object ObjectFields, cond ObjectFields) bool {
	// TODO
	return true
}

func (s *Subscription) isScalarValue(val interface{}) bool {
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
