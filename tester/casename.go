package tester

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

func CaseName(obj ...any) string {
	var s string
	for i, o := range obj {
		if i == 0 {
			if val, ok := o.(int); ok {
				s += fmt.Sprintf("%02d", val) + " "
				continue
			}
		}
		s += toString(o) + " "
	}
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`[_\s]+`).ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "/", "%")
	if len(s) > 64 {
		return s[:61] + "..."
	}
	return s
}

func toString(obj any) string {
	v := reflect.ValueOf(obj)
	var sb strings.Builder
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		// Sort the map keys to ensure consistent order
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprintf("%v", keys[i]) < fmt.Sprintf("%v", keys[j])
		})
		for _, key := range keys {
			value := v.MapIndex(key)
			sb.WriteString(fmt.Sprintf("%v ", toString(value.Interface())))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			sb.WriteString(fmt.Sprintf("%v ", toString(value.Interface())))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := v.Type().Field(i)
			if field.IsZero() || fieldType.PkgPath != "" {
				continue
			}
			sb.WriteString(toString(field.Interface()) + " ")
		}
	case reflect.Invalid:
		sb.WriteString("")
	default:
		sb.WriteString(fmt.Sprintf("%v", obj))
	}
	return strings.Trim(sb.String(), " ")
}
