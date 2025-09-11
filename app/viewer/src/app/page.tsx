export const dynamic = "force-dynamic";

import { getKV, listDummy } from "@/lib/kv";

export default async function Home() {
  const kv = await getKV();
  if (!kv) {
    console.error(
      "Cloudflare KV binding 'KV' is not available. Update wrangler.jsonc kv_namespaces and use `npm run preview`."
    );
    return (
      <div className="font-sans min-h-screen p-8 sm:p-12">
        <h1 className="text-2xl font-semibold mb-6">Cloudflare KV Dummy Viewer</h1>
        <div className="text-sm text-black/70 dark:text-white/70">データが存在しません</div>
        <p className="mt-4 text-xs text-black/60 dark:text-white/60">
          このページは prefix &quot;dummy:&quot; のキーを読み込みます。ダミーデータの投入は npm
          script を使用してください。
        </p>
      </div>
    );
  }

  let items: { key: string; value: string | null }[] = [];
  let error: string | undefined;
  try {
    items = await listDummy(kv);
  } catch (e) {
    error = e instanceof Error ? e.message : String(e);
  }

  return (
    <div className="font-sans min-h-screen p-8 sm:p-12">
      <h1 className="text-2xl font-semibold mb-6">Cloudflare KV Dummy Viewer</h1>
      {error ? (
        <div className="text-red-600">Error: {error}</div>
      ) : (
        <div className="rounded border border-black/10 dark:border-white/15">
          <div className="grid grid-cols-[200px_1fr] gap-0 text-sm">
            <div className="px-3 py-2 font-medium bg-black/[.04] dark:bg-white/[.06]">Key</div>
            <div className="px-3 py-2 font-medium bg-black/[.04] dark:bg-white/[.06]">Value</div>
            {items.map((it) => (
              <div className="contents" key={it.key}>
                <div className="px-3 py-2 border-t border-black/5 dark:border-white/10 font-mono">{it.key}</div>
                <div className="px-3 py-2 border-t border-black/5 dark:border-white/10 font-mono break-all">{it.value}</div>
              </div>
            ))}
          </div>
          {items.length === 0 && (
            <div className="px-3 py-4 text-sm text-black/70 dark:text-white/70">データが存在しません</div>
          )}
        </div>
      )}
      <p className="mt-4 text-xs text-black/60 dark:text-white/60">
        このページは prefix &quot;dummy:&quot; のキーを読み込みます。ダミーデータの投入は npm
        script を使用してください。
      </p>
    </div>
  );
}
