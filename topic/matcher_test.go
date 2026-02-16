package topic_test

import (
	"testing"

	"github.com/gorunriki/mqttc/topic"
)

func TestMatchTopic(t *testing.T) {
	test := []struct {
		filter   string
		topic    string
		expected bool
	}{
		{"sport/tennis/player1", "sport/tennis/player1", true},
		{"sport/tennis/+", "sport/tennis/player1", true},
		{"sport/tennis/#", "sport/tennis/player1/ranking", true},
		{"sport/+/player1", "sport/tennis/player1", true},
		{"sport/+/player1", "sport/soccer/player1", true},
		{"sport/+/player1", "sport/tennis/player2", false},
	}
	for _, tc := range test {
		t.Run(tc.filter+"filter to"+tc.topic, func(t *testing.T) {
			got := topic.MatchTopic(tc.filter, tc.topic)
			if got != tc.expected {
				t.Errorf("MatchTopic(%q, %q) = %v; want %v", tc.filter, tc.topic, got, tc.expected)
			}
		})
	}
}
