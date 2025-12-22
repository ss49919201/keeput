---
allowed-tools: Bash(git *)
argument-hint: [オプション: コミットメッセージ]
description: git add と git commit を自動化し、変更内容に基づいて適切なコミットメッセージを生成します
---

# Git Commit Command

## 現在の状態

**Git Status:**
!`git status`

**Unstaged Changes:**
!`git diff --stat`

**Staged Changes:**
!`git diff --staged --stat`

**Recent Commits:**
!`git log --oneline -5`

**Current Branch:**
!`git branch --show-current`

## あなたのタスク

以下の手順で git commit を作成してください:

### 1. 変更内容の分析

- 上記の git status と diff を確認
- どのファイルが変更されているか把握
- 変更の性質を理解（新機能、バグ修正、リファクタリングなど）

### 2. ファイルのステージング

- 適切なファイルを `git add` でステージング
- 不要なファイル（node_modules、.env など）は除外
- 必要に応じて `.gitignore` を追加

### 3. コミットメッセージの生成

ユーザーが引数でメッセージを指定した場合: `$ARGUMENTS`

指定がない場合、以下の形式でメッセージを生成:

```
[種類]: 簡潔な要約（50文字以内）

- 詳細な変更内容1
- 詳細な変更内容2

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

**コミットの種類:**

- `feat:` 新機能の追加
- `fix:` バグ修正
- `docs:` ドキュメントのみの変更
- `style:` フォーマットの変更
- `refactor:` リファクタリング
- `test:` テストの追加・修正
- `chore:` ビルドプロセスやツールの変更

### 4. コミット実行

**重要**: ユーザーに確認を求めず、コミットメッセージを表示してから直接実行する

heredoc を使用して複数行のコミットメッセージを作成:

```bash
git commit -m "$(cat <<'EOF'
[コミットメッセージ]
EOF
)"
```

### 5. 結果確認

- `git log -1 --stat` でコミットを確認
- `git status` で現在の状態を確認

## 注意事項

- ✅ コミット前に変更内容を確認する
- ✅ センシティブな情報は含めない
- ✅ コミットメッセージは内容を正確に表現する
- ✅ 大きすぎる変更は分割を検討する
- ✅ コミットメッセージを表示してから確認なしで直接実行する

## 例

### 例 1: 引数なしで実行

```bash
/git-commit
```

→ 変更を分析して適切なコミットメッセージを自動生成

### 例 2: メッセージを指定

```bash
/git-commit fix: ログイン時のnullエラーを修正
```

→ 指定されたメッセージでコミット
