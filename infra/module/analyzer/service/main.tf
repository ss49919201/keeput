variable "env" {
  type = string
}

variable "s3_bucket" {
  type = object({
    name = string
    arn  = string
  })
}


resource "aws_ecr_repository" "analyzer_lambda" {
  name                 = "keeput-analyzer-lambda-${var.env}"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_cloudwatch_log_group" "analyzer_lambda" {
  name              = "/aws/lambda/keeput-analyzer-${var.env}"
  retention_in_days = 14
}

data "aws_iam_policy_document" "analyzer_lambda_assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "analyzer_lambda_permissions" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["${aws_cloudwatch_log_group.analyzer_lambda.arn}:*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords"
    ]
    resources = ["*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
      "s3:PutObject"
    ]
    resources = ["${var.s3_bucket.arn}/*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket"
    ]
    resources = [var.s3_bucket.arn]
  }
}

resource "aws_iam_policy" "keeput_analyzer_lambda_policy" {
  name   = "keeput-analyzer-lambda-${var.env}"
  policy = data.aws_iam_policy_document.analyzer_lambda_permissions.json
}

resource "aws_iam_role" "analyzer_lambda" {
  name               = "keeput-analyzer-lambda-${var.env}"
  assume_role_policy = data.aws_iam_policy_document.analyzer_lambda_assume_role.json
}

resource "aws_iam_role_policy_attachment" "analyzer_lambda_policy_attachment" {
  role       = aws_iam_role.analyzer_lambda.name
  policy_arn = aws_iam_policy.keeput_analyzer_lambda_policy.arn
}

data "aws_ssm_parameter" "discord_webhook_url" {
  name = "/keeput/analyzer/discord-webhook-url"
}

data "aws_ssm_parameter" "locker_api_key_cloudflare_worker" {
  name = "/keeput/locker/api-key-cloudflare-worker"
}

data "aws_ssm_parameter" "mackerel_api_key" {
  name = "/keeput/mackerel-api-key"
}

resource "aws_lambda_function" "analyzer_lambda" {
  function_name = "keeput-analyzer-${var.env}"
  role          = aws_iam_role.analyzer_lambda.arn
  package_type  = "Image"
  memory_size   = 256
  timeout       = 30
  architectures = ["x86_64"]
  environment {
    variables = {
      DISCORD_WEBHOOK_URL              = data.aws_ssm_parameter.discord_webhook_url.value
      FEED_URL_HATENA                  = "https://ss49919201.hatenablog.com/rss"
      FEED_URL_ZENN                    = "https://zenn.dev/ss49919201/feed"
      LOCKER_API_KEY_CLOUDFLARE_WORKER = data.aws_ssm_parameter.locker_api_key_cloudflare_worker.value
      LOCKER_URL_CLOUDFLARE_WORKER     = "https://keeput-locker.ss49919201.workers.dev"
      LOG_LEVEL                        = "WARN"
      OTEL_SERVICE_NAME                = "keeput"
      S3_BUCKET_NAME                   = var.s3_bucket.name
      MACKEREL_API_KEY                 = data.aws_ssm_parameter.mackerel_api_key.value
    }
  }

  # NOTE: イメージ、関数の更新はアプリケーションのライフサイクルで行うため更新を無視する
  image_uri = "${aws_ecr_repository.analyzer_lambda.repository_url}:dummy"
  lifecycle {
    ignore_changes = [
      image_uri
    ]
  }
}
