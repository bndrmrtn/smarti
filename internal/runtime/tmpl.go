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
		// Add static content before the match
		if match[0] > lastIndex {
			result = append(result, TemplatePart{
				Static:  true,
				Content: s[lastIndex:match[0]],
			})
		}

		// Add dynamic content inside {{...}}
		result = append(result, TemplatePart{
			Static:  false,
			Content: strings.TrimSpace(s[match[2]:match[3]]),
		})

		// Update last index
		lastIndex = match[1]
	}

	// Add remaining static content after the last match
	if lastIndex < len(s) {
		result = append(result, TemplatePart{
			Static:  true,
			Content: s[lastIndex:],
		})
	}

	return result
}
