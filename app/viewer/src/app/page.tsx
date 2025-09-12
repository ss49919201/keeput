export const dynamic = "force-dynamic";

import { AnalysisReport, listAnalysisReports } from "@/lib/kv";

export default async function Home() {
  throw new Error("oops");
  const reports: { key: string; report: AnalysisReport }[] =
    await listAnalysisReports();

  return (
    <div className="font-sans min-h-screen p-8 sm:p-12">
      <h1 className="text-2xl font-semibold mb-6">
        アウトプットレポート Viewer
      </h1>
      <div className="rounded border border-black/10">
        <div className="grid grid-cols-1 sm:grid-cols-[140px_1fr_220px] lg:grid-cols-[minmax(140px,180px)_minmax(320px,1fr)_minmax(200px,260px)] gap-0 text-sm">
          <div className="px-3 py-2 font-medium bg-black/[.03]">目標達成</div>
          <div className="px-3 py-2 font-medium bg-black/[.03]">
            最新エントリ#タイトル
          </div>
          <div className="px-3 py-2 font-medium bg-black/[.03]">
            最新エントリ#公開日時
          </div>
          {reports.map((it) => {
            const le = it.report.latest_entry;
            return (
              <div className="contents" key={it.key}>
                <div className="px-3 py-2 border-t border-black/5">
                  {String(it.report.is_goal_achieved)}
                </div>
                <div className="px-3 py-2 border-t border-black/5">
                  {le ? le.Title : "-"}
                </div>
                <div className="px-3 py-2 border-t border-black/5">
                  {le ? le.PublishedAt.toISOString() : "-"}
                </div>
              </div>
            );
          })}
        </div>
        {reports.length === 0 && (
          <div className="px-3 py-4 text-sm text-black/70">
            データが存在しません
          </div>
        )}
      </div>
    </div>
  );
}
