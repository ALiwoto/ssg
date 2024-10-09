package tests_test

import (
	"log"
	"testing"

	"github.com/ALiwoto/ssg/ssg/agentUtils"
)

func TestUserAgentGeneration(t *testing.T) {
	// Test the generation of a user agent
	androidAgents := agentUtils.GetAndroidUserAgents(10)
	if len(androidAgents) != 10 {
		t.Errorf("Expected 10 android user agents, got %d", len(androidAgents))
	}

	for i, agent := range androidAgents {
		if agent.UserAgent == "" {
			t.Errorf("Expected a user agent, got an empty string")
		}

		log.Printf("User Agent %d: %s", i+1, agent.UserAgent)
	}
}
