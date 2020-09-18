package collections

import (
	"strings"
)

// GetSliceLastValueE will take a source string and returns the last value when split by the seperaror char
func GetSliceLastValueE(source string, seperator string) (string, error) {
	if len(source) > 0 && len(seperator) > 0 && strings.Contains(source, seperator) {
		tmp := strings.Split(source, seperator)
		return tmp[len(tmp)-1], nil
	}
	return "", NewSliceValueNotFoundError(source)
}

// GetSliceIndexValueE will take a source string and returns the value at the given index when split by the seperaror char
func GetSliceIndexValueE(source string, seperator string, index int) (string, error) {
	if len(source) > 0 && len(seperator) > 0 && strings.Contains(source, seperator) && index >= 0 {
		tmp := strings.Split(source, seperator)
		if index > len(tmp) {
			return "", NewSliceValueNotFoundError(source)
		}
		return tmp[index], nil
	}
	return "", NewSliceValueNotFoundError(source)
}
