package service

const (
	SEARCH_QUERY = "SEARCH_QUERY"
)

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
