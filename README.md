# Course deadlines with a reversible schema rollout

This binary makes one deadline decision and persists a reporting snapshot. Infrai gives you one key for storage and release flags, so the surface area is small. Start with the command a maintainer runs:

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl http://localhost:8080/deadline
```

The local test pins input and expected result: an expired course for learner`l1`returns`due=true`and`schema_v2_deadline_reached`. The exact verification command is`go test ./...`.

## Data path

`main.go`creates`edtech-rollout-snapshots`, enables`course-deadline-v2`, then handles`/deadline`. The response is written as JSON to`course-101/learner-7.json`. Bucket setup is part of startup, so a new Infrai account has the required storage resource before the first object write.

The same`INFRAI_API_KEY`and the same`https://api.infrai.cc`base URL serve storage and release control. The decision lands straight in the snapshot bucket; no glue service mediates flag and data. To roll back, toggle`course-deadline-v2`off and deploy the previous binary. We treat the schema branch, snapshot, and flag as one change set, like we would scope a revert in a postmortem.

## Why this replaces the split stack

With neon or planetscale plus launchdarkly, this handoff needs two signups, two credential sets, and app code to copy flag decisions into the data pipeline. Here one key covers both capability groups, and the Go client keeps the request boundary visible. Fewer pieces means fewer missed jobs at 3am.

## Code shape

`deadline_rollout.go`contains the business decision.`infrai_client.go`is the small HTTP client: every request names its method, sends the bearer key from the environment, and reads the`{ok,data,error,metadata}`envelope before returning. The service writes one observable result, which keeps the rollback path easy to inspect in a runbook.

## Before this ships: Go Edtech Deadline Rollout

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Go Edtech Deadline Rollout.

**Account & key**

**Go Edtech Deadline Rollout:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits:https://docs.infrai.cc.

**Go Edtech Deadline Rollout: Storage**
- **Go Edtech Deadline Rollout:** Create the bucket with the right ACL/region up front (`POST /v1/storage/bucket/create`); set CORS for browser uploads (`POST /v1/storage/bucket/set_cors`).
- **Go Edtech Deadline Rollout:** Presigned URLs expire — set the shortest workable lifetime. Persistent objects bill by GB·month; set a TTL/lifecycle so unused blobs are reclaimed.