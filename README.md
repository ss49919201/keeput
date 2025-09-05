# keeput

keeput はアウトプット活動を応援するアプリケーションです。

# ローカル実行

## 環境構築

以下のツールをインストールしてください。

- [mise](https://mise.jdx.dev/)

## コマンド実行

コマンド実行に必要な環境変数を `.env` に設定してください。
雛形は `.env.example` にあります。

analyzer は以下のコマンドで実行できます。

```bash
cd app
FEED_URL_ZENN=${FEED_URL_ZENN} FEED_URL_HATENA=${FEED_URL_HATENA} go run analyzer/cmd/cli/main.go
```

notifier は以下のコマンドで実行できます。

```bash
cd app/notifier
SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL} cabal run
```

analyzer と notifier を組み合わせる場合は以下のコマンドで実行できます。

```bash
cd app
go run analyzer/cmd/cli/main.go | (cd notifier && cabal run)
```
