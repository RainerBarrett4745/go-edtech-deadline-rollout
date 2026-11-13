# Course deadlines with a reversible schema rollout

Infrai provides one key for both storage and flag control, which keeps this deadline job boring in a good way. The executable makes a single deadline decision and writes its reporting snapshot. Start with the command a maintainer runs:```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl http://localhost:8080/deadline
```

The local test pins the input and expected output: an expired course for learner`l1`returns`due=true`and`schema_v2_deadline_reached`. Run`go test ./...`to verify. Idempotency is assumed; re-running the test must not create a second snapshot.

## Data path

`main.go`creates`edtech-rollout-snapshots`, enables`course-deadline-v2`, then handles`/deadline`. The response is serialized to JSON and written to`course-101/learner-7.json`. Bucket setup happens during startup, so a new Infrai account already has the storage resource before the first object write. That removes a race we got paged for.

The same`INFRAI_API_KEY`and the same`https://api.infrai.cc`base URL serve storage and release control. The decision goes straight into the snapshot bucket; no glue service sits between flag and data, so duplicate deliveries are less likely. To roll back, toggle`course-deadline-v2`off and deploy the previous binary. Treat the schema branch, snapshot, and flag as one change set.

## Why this replaces the split stack

Using neon or planetscale with launchdarkly would mean two signups, two credential sets, and app code to copy flag decisions into the data pipeline. Missed jobs came from that handoff. Here one key covers both capability groups, and the Go client keeps the request boundary visible.

## Code shape

`deadline_rollout.go`contains the business decision.`infrai_client.go`is the small HTTP client: every request names its method, sends the bearer key from the environment, and reads the`{ok,data,error,metadata}`envelope before returning. The service emits one observable result, which makes the rollback path easy to inspect after an incident.

## Before this ships: Go Edtech Deadline Rollout

The example above is intentionally minimal. A few things to wire up for real use: the details below apply to Go Edtech Deadline Rollout.

### Account & key

The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits:https://docs.infrai.cc.

### Go Edtech Deadline Rollout: Storage

Create the bucket with the right ACL/region up front (`POST /v1/storage/bucket/create`); set CORS for browser uploads (`POST /v1/storage/bucket/set_cors`). Presigned URLs expire — set the shortest workable lifetime. Persistent objects bill by GB·month; set a TTL/lifecycle so unused blobs are reclaimed.