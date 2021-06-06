package main

// Check if the string exists in a slice
func containString(sl []string, value string) bool {
	for _, str := range sl {
		if value == str {
			return true
		}
	}

	return false
}