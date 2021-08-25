package handler

import (
	"strconv"
	"time"
)

const (
	SEARCH_QUERY = "SEARCH_QUERY"
)

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func unixStringToRFC3339String(date string) (string, error) {
	sinceDateUInt, err := strconv.ParseUint(string(date), 10, 64)
	if err != nil {
		return "", err
	}
	sinceDateString := time.Unix(int64(sinceDateUInt), int64(0)).UTC().Format(time.RFC3339)
	return sinceDateString, nil
}
