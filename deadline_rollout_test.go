package main

import (
	"testing"
	"time"
)

func TestDecideDeadlineFlaggedSchema(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	course := Course{ID: "c1", Deadline: now.Add(-time.Hour)}
	got := decideDeadline(course, Learner{ID: "l1"}, now, true)
	if !got.Due || got.Reason != "schema_v2_deadline_reached" {
		t.Fatalf("unexpected decision: %+v", got)
	}
}

func TestDecideDeadlineLegacyPending(t *testing.T) {
	now := time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)
	course := Course{ID: "c1", Deadline: now.Add(time.Hour)}
	got := decideDeadline(course, Learner{ID: "l1"}, now, false)
	if got.Due || got.Reason != "deadline_pending" {
		t.Fatalf("unexpected decision: %+v", got)
	}
}
