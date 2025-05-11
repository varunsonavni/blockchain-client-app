provider "aws" {
  region = var.region
}

# Create ECR repository for blockchain client images
resource "aws_ecr_repository" "blockchain_client" {
  name                 = "${var.environment}-blockchain-client"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name        = "${var.environment}-blockchain-client-repo"
    Environment = var.environment
  }
}

# ECS Fargate deployment for blockchain client
module "ecs" {
  source = "../../modules/ecs"

  environment        = var.environment
  region             = var.region
  vpc_id             = data.terraform_remote_state.rs.outputs.vpc_id
  public_subnet_ids  = data.terraform_remote_state.rs.outputs.public_subnets
  private_subnet_ids = data.terraform_remote_state.rs.outputs.private_subnets

  # Container configuration
  container_image    = aws_ecr_repository.blockchain_client.repository_url
  container_version  = var.container_version
  blockchain_rpc_url = var.blockchain_rpc_url

  # Task and service configuration
  task_cpu              = var.task_cpu
  task_memory           = var.task_memory
  service_desired_count = var.service_desired_count
  service_min_capacity  = var.service_min_capacity
  service_max_capacity  = var.service_max_capacity
  log_retention_days    = var.log_retention_days
} 