import { getCloudflareContext } from "@opennextjs/cloudflare";

export type KvItem = { key: string; value: string | null };

export async function getKV(): Promise<KVNamespace | undefined> {
  // Use async mode to comply with Next.js restrictions
  const ctx = await getCloudflareContext({ async: true });
  const env = (ctx?.env ?? {}) as Partial<CloudflareEnv>;
  return env.KV as unknown as KVNamespace | undefined;
}

// Seeding is handled via npm scripts now.

export async function listDummy(kv?: KVNamespace): Promise<KvItem[]> {
  const prefix = "dummy:";
  if (!kv) throw new Error("KV binding is not available");
  const res = await kv.list({ prefix });
  const items: KvItem[] = [];
  for (const k of res.keys) {
    const v = await kv.get(k.name);
    items.push({ key: k.name.replace(prefix, ""), value: v });
  }
  return items;
}
