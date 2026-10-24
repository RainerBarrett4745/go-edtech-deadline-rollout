package main

import "time"

type Course struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Deadline time.Time `json:"deadline"`
}

type Learner struct {
	ID string `json:"id"`
}

type DeadlineDecision struct {
	CourseID  string `json:"course_id"`
	LearnerID string `json:"learner_id"`
	Due       bool   `json:"due"`
	Reason    string `json:"reason"`
}

func decideDeadline(course Course, learner Learner, now time.Time, newSchema bool) DeadlineDecision {
	due := !now.Before(course.Deadline)
	reason := "deadline_pending"
	if due {
		reason = "deadline_reached"
	}
	if newSchema {
		reason = "schema_v2_" + reason
	}
	return DeadlineDecision{CourseID: course.ID, LearnerID: learner.ID, Due: due, Reason: reason}
}
