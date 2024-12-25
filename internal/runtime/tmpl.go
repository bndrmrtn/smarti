package runtime

import (
	"regexp"
	"strings"
)

type TemplatePart struct {
	Static  bool
	Content string
}

func parseTemplate(s string) []TemplatePart {
	var result []TemplatePart
	re := regexp.MustCompile(`\{\{(.*?)}}`)
	matches := re.FindAllStringSubmatchIndex(s, -1)

	lastIndex := 0
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		if match[0] > lastIndex {
			result = append(result, TemplatePart{
				Static:  true,
				Content: s[lastIndex:match[0]],
			})
		}

		result = append(result, TemplatePart{
			Static:  false,
			Content: strings.TrimSpace(s[match[2]:match[3]]),
		})

		lastIndex = match[1]
	}

	if lastIndex < len(s) {
		result = append(result, TemplatePart{
			Static:  true,
			Content: s[lastIndex:],
		})
	}

	return result
}
