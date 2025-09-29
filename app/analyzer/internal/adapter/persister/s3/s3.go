package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"path"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ss49919201/keeput/app/analyzer/internal/appctx"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
	"github.com/ss49919201/keeput/app/analyzer/internal/model"
	"github.com/ss49919201/keeput/app/analyzer/internal/port/persister"
)

var s3ClientOnce sync.Once
var s3Client *s3.Client

func initS3Client(config aws.Config) {
	s3ClientOnce.Do(func() {
		s3Client = s3.NewFromConfig(config)
	})
}

func NewPersistAnalysisReport(config aws.Config) persister.PersistAnalysisReport {
	initS3Client(config)
	return persistAnalysisReport
}

func persistAnalysisReport(ctx context.Context, report *model.AnalysisReport) error {
	b, err := json.Marshal(report)
	if err != nil {
		return err
	}

	bucket := config.S3BucketName()
	now := appctx.GetNowOr(ctx, time.Now())
	key := path.Join("analysis_report", now.Format("2006/01/02/15/04/05"), "data.json")

	if _, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(b),
		ContentType: aws.String("application/json"),
	}); err != nil {
		return err
	}

	return s3.NewObjectExistsWaiter(s3Client).Wait(
		ctx,
		&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
		time.Minute,
	)
}
