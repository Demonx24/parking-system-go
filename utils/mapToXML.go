package utils

import "encoding/xml"

func MapToXML(params map[string]string) string {
	type Entry struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}

	xmlData := struct {
		XMLName xml.Name `xml:"xml"`
		Entries []Entry
	}{}

	for k, v := range params {
		xmlData.Entries = append(xmlData.Entries, Entry{
			XMLName: xml.Name{Local: k},
			Value:   v,
		})
	}

	output, _ := xml.Marshal(xmlData)
	return string(output)
}
