package db

import (
	"client"
	"errors"
	"logger"
	"math"
	"reflect"
	. "types"
)

type Subscription struct {
	Key     string
	Query   ObjectFields
	Clients []*client.Client
}

// Tells if subscription query valid
func (s *Subscription) Validate() error {
	for key, condition := range s.Query {
		if err := s.validatePair(key, condition); err != nil {
			return err
		}
	}
	return nil
}

// Recursively validates pairs of key/condition
func (s *Subscription) validatePair(key string, condition interface{}) error {
	switch key {
	case "OR", "AND", "XOR":
		if _, ok := condition.(map[string]interface{}); !ok {
			return errors.New("Key OR/AND/XOR must have list of conditions")
		}
		if len(condition.(map[string]interface{})) < 2 {
			return errors.New("Key OR/AND/XOR must have more than one condition")
		}
	case "NOT":
		if _, ok := condition.(map[string]interface{}); !ok {
			return errors.New("Key NOT must have list of single condition")
		}
		if len(condition.(map[string]interface{})) != 1 {
			return errors.New("Key NOT must have only one condition")
		}
	}
	if _, ok := condition.(map[string]interface{}); ok {
		for key2, condition2 := range condition.(map[string]interface{}) {
			if err := s.validatePair(key2, condition2); err != nil {
				return err
			}
		}
	}
	return nil
}

// Match the object against subscription's Query
func (s *Subscription) Match(object ObjectFields) bool {
	// No conditions, pass all the object updates
	if len(s.Query) == 0 {
		return true
	}
	if len(s.Query) == 1 {
		for key, condition := range s.Query {
			return s.match(key, condition, object)
		}
	}
	return s.match("AND", s.Query, object)
}

func (s *Subscription) match(key string, condition interface{}, object ObjectFields) bool {
	if c, ok := condition.(map[string]interface{}); ok {
		switch key {
		case "OR": // At least one condition is true
			for key2, condition2 := range c {
				outcome := s.match(key2, condition2, object)
				if outcome {
					return true
				}
			}
			return false
		case "AND": // Every condition is true
			for key2, condition2 := range c {
				outcome := s.match(key2, condition2, object)
				if !outcome {
					return false
				}
			}
			return true
		case "XOR": // Must be odd number of true conditions
			result := false
			for key2, condition2 := range c {
				if outcome := s.match(key2, condition2, object); outcome {
					result = !result
				}
			}
			return result
		case "NOT": // Negate single result
			for key2, condition2 := range c {
				return !s.match(key2, condition2, object)
			}
		default:
			logger.ErrorLog.Printf("Key %s not supported", key)
		}
	} else if s.isScalarValue(condition) { // Plain value, perform comparison
		if reflect.TypeOf(condition) != reflect.TypeOf(object[key]) {
			f1, err1 := s.getFloat(condition)
			f2, err2 := s.getFloat(object[key])
			if err1 == nil && err2 == nil {
				return f1 == f2
			} else {
				logger.ErrorLog.Printf(
					"Type mismatch for property %s: wanted %s, got %s",
					key,
					reflect.TypeOf(condition),
					reflect.TypeOf(object[key]),
				)
				return false
			}
		}
		return condition == object[key]
	}
	return false
}

func (s *Subscription) getFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int:
		return float64(v), nil
	default:
		return math.NaN(), errors.New("Not a number")
	}
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
