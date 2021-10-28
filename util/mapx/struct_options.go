package mapx

import (
	"strings"
)

type structOption interface {
	getOptionName() string
}

type structKeyFieldOption struct {
	mappings [][2]string
}

type structKeyFieldEntry struct {
	key       string
	fieldName string
}

func (opt *structKeyFieldOption) getOptionName() string {
	return "StructKeyFieldOption"
}

func (opt *structKeyFieldOption) getEntries() []structKeyFieldEntry {
	entries := make([]structKeyFieldEntry, 0)

	for _, parts := range opt.mappings {
		p1 := strings.TrimSpace(parts[0])
		p2 := strings.TrimSpace(parts[1])

		if p1 == "" || p2 == "" || p1 == p2 {
			continue
		}

		entries = append(entries, structKeyFieldEntry{key: p1, fieldName: p2})
	}

	return entries
}

type structKeysOption struct {
	mode string
	keys []string
}

func (opt *structKeysOption) getOptionName() string {
	return "StructKeysOption"
}
