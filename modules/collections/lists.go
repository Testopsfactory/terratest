package collections

// Remove all the items in list2 from list1
func Subtract(list1 []string, list2 []string) []string {
	out := []string{}

	for _, item := range list1 {
		if !ListContains(item, list2) {
			out = append(out, item)
		}
	}

	return out
}

// Return true if the given list of strings (haystack) contains the given string (needle)
func ListContains(needle string, haystack []string) bool {
	for _, str := range haystack {
		if needle == str {
			return true
		}
	}

	return false
}
