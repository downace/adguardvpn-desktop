package common

import (
	"regexp"
	"strings"
)

type tableField struct {
	length int
	name   string
}

var headerRegexp = regexp.MustCompile(`(\w+)\s+`)

func ParseTable[T any](table string, convert func(row map[string]string) T) []T {
	fields := make([]tableField, 0)
	result := make([]T, 0)
	isHeader := true
	for line := range strings.Lines(table) {
		if isHeader {
			for _, match := range headerRegexp.FindAllStringSubmatch(line, -1) {
				fields = append(fields, tableField{len(match[0]), match[1]})
			}
		} else if strings.TrimSpace(line) == "" {
			break
		} else {
			row := make(map[string]string)
			rest := line
			for _, field := range fields {
				l := min(field.length, len(rest))
				value := strings.TrimSpace(rest[:l])
				rest = rest[l:]

				row[field.name] = value
			}
			result = append(result, convert(row))
		}

		isHeader = false
	}
	return result
}
