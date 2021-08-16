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
	fromDateUInt, err := strconv.ParseUint(string(date), 10, 64)
	if err != nil {
		return "", err
	}
	fromDateString := time.Unix(int64(fromDateUInt), int64(0)).UTC().Format(time.RFC3339)
	return fromDateString, nil
}
