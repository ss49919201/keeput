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

module "service" {
  source = "../../../../module/analyzer/service"
  env    = "dev"
}