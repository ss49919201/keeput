## KV AnalysisReport Viewer

- UI: Home page renders a table of AnalysisReport JSON stored in KV (prefix `analysis_report:`).
- Schema matches Go side definitions:
  - `is_goal_achieved: boolean`
  - `latest_entry: { Title: string, Body: string, PublishedAt: string } | null`
  - Entry field names follow Go JSON (Title/Body/PublishedAt)
- Seeding is done via npm scripts, not in runtime.

### Run locally on Cloudflare (preview)

```
npm run preview
```

Wrangler binds a KV namespace as `KV` (see `wrangler.jsonc`).
