package topic

import "strings"

func MatchTopic(filter, topic string) bool {
	if filter == topic {
		return true
	}

	return matchesWithWildcards(filter, topic)
}

func matchesWithWildcards(filter, topic string) bool {
	filterParts := strings.Split(filter, "/")
	topicParts := strings.Split(topic, "/")

	for i := 0; i < len(filterParts) && i < len(topicParts); i++ {
		if filterParts[i] == "#" {
			return true
		}
		if filterParts[i] == "+" {
			continue
		}
		if filterParts[i] != topicParts[i] {
			return false
		}
	}
	return len(filterParts) == len(topicParts)
}
