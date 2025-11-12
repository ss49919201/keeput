terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.20.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

module "storage" {
  source = "../../../module/analyzer/storage"
  env    = "prod"
}

module "service" {
  source    = "../../../module/analyzer/service"
  env       = "prod"
  s3_bucket = module.storage.s3_bucket
}