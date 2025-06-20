package resource_virtualmachine

import (
	"fmt"
	"sort"
	"strings"
)

func BuildMetadataString(tags map[string]interface{}, description, notes string) string {
	var parts []string

	if len(tags) > 0 {
		// Make tag output stable (for testing and diffs)
		var tagPairs []string
		for k, v := range tags {
			tagPairs = append(tagPairs, fmt.Sprintf("%s: %v", k, v))
		}
		sort.Strings(tagPairs)
		parts = append(parts, fmt.Sprintf("{Tags: %s}", strings.Join(tagPairs, ", ")))
	}

	if description != "" {
		parts = append(parts, fmt.Sprintf("{description: %s}", description))
	}

	if notes != "" {
		parts = append(parts, fmt.Sprintf("{notes: %s}", notes))
	}

	return strings.Join(parts, "")
}

func ParseMetadataString(metadata string) (map[string]interface{}, string, string) {
	tags := map[string]interface{}{}
	var description, notes string

	for metadata != "" {
		start := strings.Index(metadata, "{")
		end := strings.Index(metadata, "}")
		if start == -1 || end == -1 || end < start {
			break
		}

		section := metadata[start+1 : end]
		metadata = metadata[end+1:]

		switch {
		case strings.HasPrefix(section, "Tags:"):
			tagString := strings.TrimSpace(strings.TrimPrefix(section, "Tags:"))
			pairs := strings.Split(tagString, ",")
			for _, pair := range pairs {
				if kv := strings.SplitN(strings.TrimSpace(pair), ":", 2); len(kv) == 2 {
					key := strings.TrimSpace(kv[0])
					value := strings.TrimSpace(kv[1])
					tags[key] = value
				}
			}
		case strings.HasPrefix(section, "description:"):
			description = strings.TrimSpace(strings.TrimPrefix(section, "description:"))
		case strings.HasPrefix(section, "notes:"):
			notes = strings.TrimSpace(strings.TrimPrefix(section, "notes:"))
		}
	}

	return tags, description, notes
}
