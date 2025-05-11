terraform {
  backend "s3" {
    bucket = "terraform-state-dev-8008"
    key    = "dev/ecs/terraform.tfstate"
    region = "us-east-1"
  }
}

