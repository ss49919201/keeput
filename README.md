# keeput

keeput はアウトプット活動を応援するアプリケーションです。

# 実行方法

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
FEED_URL_ZENN=${FEED_URL_ZENN} FEED_URL_HATENA=${FEED_URL_HATENA} SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL} go run analyzer/cmd/cli/main.go | (cd notifier && cabal run)
```
