package util

// AppendIfDefined appends e to arr if e is defined.
func AppendIfDefined(arr []string, e string) []string {
	if len(e) == 0 {
		return arr
	}

	return append(arr, e)
}
