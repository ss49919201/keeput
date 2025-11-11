# keeput

keeput はアウトプット活動を応援するアプリケーションです。

# locker

## 概要

複数の analyzer インスタンスが同じユースケースを同時に実行しないよう排他制御を行うための分散ロックサービスです。

Durable Object の特性を活用し、10 秒間リクエストがない場合に自動的にロックが解放されます。

## 注意事項

- ロック ID ごとに個別の Durable Object インスタンスが作成されます。
- 自動 TTL は約 10 秒ですが、正確な時間は保証されません。
- 明示的に release を呼ぶことを推奨します。

## ローカル実行

必要な環境変数を `./app/locker/.env` に設定してください。
雛形は `./app/locker/.env.example` にあります。

```bash
cd app/locker
npm run dev
```

# analyzer

## 概要

アウトプット状況を分析するサービスです。
分析結果はチャットサービスに通知されます。

# ローカル実行

実動作環境は AWS Lambda を想定しているため、ローカル実行には CLI 用のインターフェースを使用します。

必要な環境変数を `./app/.env` に設定してください。
雛形は `./app.env.example` にあります。

```bash
cd app
ENV=local go run analyzer/cmd/cli/main.go
```
