package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func main() {
	duration, _ := ToDurationE("60")

	fmt.Println(duration * time.Second)
}

// ToDurationE method
func ToDurationE(i interface{}) (d time.Duration, err error) {
	i = indirect(i)

	switch s := i.(type) {
	case time.Duration:
		return s, nil
	case int64, int32, int16, int8, int:
		d = time.Duration(ToInt64(s))
		return
		/*	case float32, float64:
			d = time.Duration(ToFloat64(s))
			return*/
	case string:
		if strings.ContainsAny(s, "nsuµmh") {
			d, err = time.ParseDuration(s)
		} else {
			d, err = time.ParseDuration(s + "ns")
		}
		return
	default:
		err = fmt.Errorf("unable to cast %#v to duration", i)
		return
	}
}

// ToInt64 is a method
func ToInt64(i interface{}) int64 {
	v, _ := ToInt64E(i)
	return v
}

// ToInt64E casts an empty interface to an int64.
func ToInt64E(i interface{}) (int64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case int64:
		return s, nil
	case int:
		return int64(s), nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
		/*	case string:
			v, err := strconv.ParseInt(s, 0, 0)
			if err == nil {
				return v, nil
			}
			return 0, fmt.Errorf("Unable to Cast %#v to int64", i)*/
	case float64:
		return int64(s), nil
	case bool:
		if bool(s) {
			return int64(1), nil
		}
		return int64(0), nil
	case nil:
		return int64(0), nil
	default:
		return int64(0), fmt.Errorf("Unable to Cast %#v to int64", i)
	}
}

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
