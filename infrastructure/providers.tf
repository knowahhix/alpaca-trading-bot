terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    bucket = "alpaca-infra"
    key    = "terraform.tfstate"
    region = "us-east-1"
  }
}

# Configure the AWS Provider
variable "AWS_REGION" {
  description = "region"
}

variable "AWS_SECRET_ACCESS_KEY" {
  description = "secret key"
}

variable "AWS_ACCESS_KEY_ID" {
  description = "access key"
}

provider "aws" {
  region = "us-east-1"
}

