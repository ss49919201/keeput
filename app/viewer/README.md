## KV Dummy Viewer

- UI: Home page renders a simple table of KV key/value pairs (prefix `dummy:`).
- Seeding is done via npm scripts, not in runtime.

### Run locally on Cloudflare (preview)

```
npm run preview
```

Wrangler binds a KV namespace as `KV` (see `wrangler.jsonc`).
