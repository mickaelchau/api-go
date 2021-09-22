terraform {
  required_version = "1.0.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
    local = {
      source = "hashicorp/local"
      version = "2.1.0"
    }
  }
}

provider "aws" {
  region                  = "eu-west-3" // file ?
  shared_credentials_file = "/home/micka/.aws/credentials"
}
