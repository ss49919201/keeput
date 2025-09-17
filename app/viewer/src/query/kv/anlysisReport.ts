import { getCloudflareContext } from "@opennextjs/cloudflare";
import {
  boolean,
  InferOutput,
  object,
  parse,
  pipe,
  string,
  transform,
  union,
  null_ as vNull,
} from "valibot";

export type KvItem = { key: string; value: string | null };

export async function getKV(): Promise<KVNamespace | undefined> {
  const ctx = await getCloudflareContext({ async: true });
  const env = (ctx?.env ?? {}) as Partial<CloudflareEnv>;
  return env.KV as unknown as KVNamespace | undefined;
}

const dateFromRFC3339 = pipe(
  string(),
  transform((s) => new Date(s))
);

// Incoming shape stored in KV (follows Go JSON tags / Title-case & snake_case)
const entryInSchema = object({
  Title: string(),
  Body: string(),
  PublishedAt: dateFromRFC3339,
});

// Exposed shape (camelCase)
const analysisReportSchema = pipe(
  object({
    is_goal_achieved: boolean(),
    latest_entry: union([entryInSchema, vNull()]),
  }),
  transform((r) => ({
    isGoalAchieved: r.is_goal_achieved,
    latestEntry: r.latest_entry
      ? {
          title: r.latest_entry.Title,
          body: r.latest_entry.Body,
          publishedAt: r.latest_entry.PublishedAt,
        }
      : null,
  }))
);

export type AnalysisReport = InferOutput<typeof analysisReportSchema>;

export type ParsedReport = {
  key: string;
  report: AnalysisReport;
};

export async function listAnalysisReports(): Promise<ParsedReport[]> {
  const kv = (await getCloudflareContext({ async: true })).env.KV;
  const prefix = "analysis_report";
  const res = await kv.list({ prefix });
  const items: ParsedReport[] = [];
  for (const k of res.keys) {
    const raw = await kv.get(k.name);
    if (!raw) continue;
    try {
      const parsed = parse(analysisReportSchema, JSON.parse(raw));
      items.push({ key: k.name.replace(prefix, ""), report: parsed });
    } catch {
      // ignore: invalid schema or JSON
    }
  }
  return items;
}

