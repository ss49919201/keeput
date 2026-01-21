package persister

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	persisterport "github.com/ss49919201/keeput/app/analyzer/internal/port/persister"
)

func NewPersistEntry(s3client *s3.Client) persisterport.PersistEntry {
	return func(ctx context.Context, entry *model.Entry) error {
		return persistEntry(ctx, s3client, entry)
	}
}

// TODO
func persistEntry(ctx context.Context, s3client *s3.Client, entry *model.Entry) error {
	// S3 から sqlite ファイル取得
	// sqlite ファイルをエファメラルストレージに保存
	// sql.DB 初期化
	// entry 挿入
	// S3 にアップロード
	return nil
}
