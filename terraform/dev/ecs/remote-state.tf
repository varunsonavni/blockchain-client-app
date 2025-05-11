data "terraform_remote_state" "rs" {
  backend = "s3"

  config = {
    bucket = "terraform-state-dev-8008"
    key    = "dev/vpc/terraform.tfstate"
    region = "us-east-1"
  }
}