"use client";

import Link from "next/link";

export default function Error() {
  return (
    <div className="font-sans min-h-screen p-8 sm:p-12">
      <div className="flex items-center gap-3 mb-4">
        <div className="text-3xl select-none animate-bounce" aria-hidden>
          🫠
        </div>
        <h1 className="text-2xl font-semibold">オイオイ！エラーだぞ</h1>
      </div>

      <div className="mt-3">
        <div className="inline-block rounded-md border border-black/10 px-4 py-2 bg-black/5 text-sm font-mono">
          ¯\\_(ツ)_/¯ なんか調子悪いっぽい
        </div>
      </div>

      <div className="mt-5 flex flex-wrap items-center gap-3">
        <Link
          href="/"
          className="px-3 py-1.5 rounded border border-sky-400/50 bg-sky-50 text-sky-900 text-sm hover:bg-sky-100/80 dark:hover:bg-sky-900/50"
        >
          トップへ戻る
        </Link>
      </div>
    </div>
  );
}
