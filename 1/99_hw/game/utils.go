package main

// Check if the string exists in a slice
func containString(sl []string, value string) (bool, int) {
	for index, str := range sl {
		if value == str {
			return true, index
		}
	}

	return false, -1
}

func removeStringFromSlice(sl []string, index int) []string {
	ret := make([]string, 0)
	if len(sl) <= 1 {
		return ret
	}
	ret = append(ret, sl[:index]...)
	return append(ret, sl[index+1:]...)
}
