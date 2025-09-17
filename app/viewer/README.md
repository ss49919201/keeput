## KV AnalysisReport Viewer

- UI: Home page renders a table of AnalysisReport JSON stored in KV (prefix `analysis_report:`).
- Exposed schema (Viewer output) uses camelCase:
  - `isGoalAchieved: boolean`
  - `latestEntry: { title: string, body: string, publishedAt: Date } | null`
- Stored schema (in KV) matches Go JSON:
  - `is_goal_achieved: boolean`
  - `latest_entry: { Title: string, Body: string, PublishedAt: string(RFC3339) } | null`
  - Viewer transforms stored fields to camelCase on read.
- Seeding is done via npm scripts, not in runtime.

### Run locally on Cloudflare (preview)

```
npm run preview
```

Wrangler binds a KV namespace as `KV` (see `wrangler.jsonc`).
