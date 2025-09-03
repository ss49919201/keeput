# Architecture Decision Records (ADR)

このディレクトリは keeput システムのアーキテクチャ決定記録を管理します。

## ADR について

Architecture Decision Records (ADR) は、アーキテクチャに関する重要な決定を記録する軽量な文書です。

### ADR の構成

- **ステータス**: 提案中、承認済み、廃止済み等
- **コンテキスト**: 決定が必要になった背景・状況
- **決定事項**: 何を決定したか
- **理由**: なぜその決定をしたか
- **影響**: 決定による影響・トレードオフ

### 新しい ADR の作成

1. `adr/` ディレクトリに `00X-title.md` 形式でファイルを作成
2. 上記テンプレートに沿って内容を記述
3. この README の一覧表を更新

### 参考資料

- [ADR GitHub](https://adr.github.io/)
- [Architecture decision record](https://en.wikipedia.org/wiki/Architectural_decision)
