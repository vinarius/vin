package watch

import "strings"

func getTopicsFromTopicFlag(topicRaw string) []string {
	return strings.Split(topicRaw, ",")
}
