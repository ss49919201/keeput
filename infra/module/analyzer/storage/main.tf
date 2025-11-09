variable "env" {
  type = string
}

resource "aws_s3_bucket" "analyzer" {
  bucket = "keeput-analyzer-${var.env}"
}

resource "aws_s3_bucket_public_access_block" "analyzer" {
  bucket = aws_s3_bucket.analyzer.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "analyzer" {
  bucket = aws_s3_bucket.analyzer.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "analyzer" {
  bucket = aws_s3_bucket.analyzer.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket" "analyzer_log" {
  bucket = "keeput-analyzer-log-${var.env}"
}

resource "aws_s3_bucket_logging" "example" {
  bucket = aws_s3_bucket.analyzer.id
  target_bucket = aws_s3_bucket.analyzer_log.id
  target_prefix = "log/"
}

output "bucket_name" {
  value       = aws_s3_bucket.analyzer.id
}

output "bucket_arn" {
  value       = aws_s3_bucket.analyzer.arn
}
