package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const snapshotBucket = "edtech-rollout-snapshots"
const schemaFlag = "course-deadline-v2"

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.CreateBucket(snapshotBucket); err != nil {
		log.Fatal(err)
	}
	if err := client.ToggleFlag(schemaFlag, true, 1); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/deadline", func(w http.ResponseWriter, r *http.Request) {
		course := Course{ID: "course-101", Title: "Data Pipelines", Deadline: time.Now().Add(24 * time.Hour)}
		learner := Learner{ID: "learner-7"}
		decision := decideDeadline(course, learner, time.Now(), true)
		payload, _ := json.Marshal(decision)
		if err := client.WriteSnapshot(snapshotBucket, "course-101/learner-7.json", payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	})
	fmt.Println("deadline service listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
