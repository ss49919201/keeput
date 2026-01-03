terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.20.0"
    }
  }
  backend "s3" {
    bucket = "keeput-tf-state"
    key    = "prod/analyzer/service/terraform.tfstate"
    region = "ap-northeast-1"
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

module "storage" {
  source          = "../../../module/analyzer/storage"
  env             = "prod"
}

module "service" {
  source    = "../../../module/analyzer/service"
  env       = "prod"
  scheduler_state = "ENABLED"
  s3_bucket = module.storage.s3_bucket
}