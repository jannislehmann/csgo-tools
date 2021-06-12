package util

// Reverse takes a string and reverses it.
func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
