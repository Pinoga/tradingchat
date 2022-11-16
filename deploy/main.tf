provider "aws" {
  region = "us-east-1"
}

terraform {
  required_version = "1.3.4"
}

resource "aws_ecr_repository" "services" {
  name                 = "services"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}
