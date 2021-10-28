package mapx

import (
	"encoding/xml"
	"io"
	"strings"
)

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type stringMap map[string]string

func (m *stringMap) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	*m = stringMap{}

	for {
		var entry xmlMapEntry
		err := decoder.Decode(&entry)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[entry.XMLName.Local] = entry.Value
	}

	return nil
}

type xmlOption interface {
	getOptionName() string
}

type xmlKeyTagOption struct {
	mappings [][2]string
}

type xmlKeyTagEntry struct {
	key string
	tag string
}

func (opt *xmlKeyTagOption) getOptionName() string {
	return "XmlKeyTagOption"
}

func (opt *xmlKeyTagOption) getEntries() []xmlKeyTagEntry {
	entries := make([]xmlKeyTagEntry, 0)

	for _, parts := range opt.mappings {
		p1 := strings.TrimSpace(parts[0])
		p2 := strings.TrimSpace(parts[1])

		if p1 == "" || p2 == "" || p1 == p2 {
			continue
		}

		entries = append(entries, xmlKeyTagEntry{key: p1, tag: p2})
	}

	return entries
}

type xmlKeysOption struct {
	mode string
	keys []string
}

func (opt *xmlKeysOption) getOptionName() string {
	return "XmlKeysOption"
}

type xmlCDataOption struct {
	keys []string
}

func (opt *xmlCDataOption) getOptionName() string {
	return "XmlCDataOption"
}
