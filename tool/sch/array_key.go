package sch

import "fmt"

func AssertKey[K comparable](queryName string, queryFunc func(any) K) (string, func(_ map[string]any, field any) error) {
	return "[*]." + queryName + " is unique", func(_ map[string]any, field any) error {
		items := field.([]any)
		knownValues := make(map[K]bool, len(items))
		for _, item := range items {
			key := queryFunc(item)
			if knownValues[key] {
				return fmt.Errorf("Multiple key: %v", key)
			}
			knownValues[key] = true
		}
		return nil
	}
}

func AssertOnlyOneTrue(queryName string, queryFunc func(any) bool) (string, func(_ map[string]any, field any) error) {
	return "exact one [$]." + queryName + " is true", func(_ map[string]any, field any) error {
		items := field.([]any)
		count := 0
		for _, item := range items {
			if queryFunc(item) {
				count++
			}
		}
		if count == 0 {
			return fmt.Errorf("Any [$].%s is true", queryName)
		} else if count > 1 {
			return fmt.Errorf("Multiple [$].%s is true", queryName)
		}
		return nil
	}
}
