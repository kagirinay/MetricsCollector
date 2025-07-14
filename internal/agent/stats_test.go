package agent

import "testing"

func TestStatsCollectorPoll(t *testing.T) {
	c := NewStatsCollector()
	c.poll()
	if c.PollCount != 1 {
		t.Fatalf("PollCount expected 1, got %d", c.PollCount)
	}
	if _, ok := c.gauges["RandomValue"]; !ok {
		t.Fatalf("RandomValue gauge not present")
	}
}
