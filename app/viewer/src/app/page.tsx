export const dynamic = "force-dynamic";

import { AnalysisReport, listAnalysisReports } from "@/query/kv/anlysisReport";

export default async function Home() {
  const reports: { key: string; report: AnalysisReport }[] =
    await listAnalysisReports();

  return (
    <main className="min-h-screen p-6 sm:p-10 font-sans">
      <h1 className="text-xl font-semibold mb-4">アウトプットレポート Viewer</h1>
      {reports.length === 0 ? (
        <p className="text-gray-600">データが存在しません</p>
      ) : (
        <div className="rounded-md border border-gray-200 overflow-hidden">
          <table className="w-full border-collapse text-sm">
            <thead>
              <tr>
                <th className="text-left bg-gray-50 border border-gray-200 p-2 sm:p-3">目標達成</th>
                <th className="text-left bg-gray-50 border border-gray-200 p-2 sm:p-3">最新エントリ#タイトル</th>
                <th className="text-left bg-gray-50 border border-gray-200 p-2 sm:p-3">最新エントリ#公開日時</th>
              </tr>
            </thead>
            <tbody>
              {reports.map((it) => {
                const le = it.report.latestEntry;
                return (
                  <tr key={it.key} className="odd:bg-white even:bg-gray-50">
                    <td className="border border-gray-200 p-2 sm:p-3">{String(it.report.isGoalAchieved)}</td>
                    <td className="border border-gray-200 p-2 sm:p-3">{le ? le.title : "-"}</td>
                    <td className="border border-gray-200 p-2 sm:p-3">{le ? le.publishedAt.toISOString() : "-"}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </main>
  );
}
