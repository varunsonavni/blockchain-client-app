variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs"
  type        = list(string)
}

variable "container_image" {
  description = "Docker image repository for the blockchain client"
  type        = string
}

variable "container_version" {
  description = "Docker image version/tag for the blockchain client"
  type        = string
  default     = "latest"
}

variable "blockchain_rpc_url" {
  description = "Blockchain RPC URL (e.g., Polygon RPC endpoint)"
  type        = string
  default     = "https://polygon-rpc.com/"
}

variable "task_cpu" {
  description = "CPU units for the ECS task (1024 = 1 vCPU)"
  type        = number
  default     = 256
}

variable "task_memory" {
  description = "Memory for the ECS task in MiB"
  type        = number
  default     = 512
}

variable "service_desired_count" {
  description = "Number of instances of the task to run"
  type        = number
  default     = 2
}

variable "service_min_capacity" {
  description = "Minimum number of instances of the task for auto scaling"
  type        = number
  default     = 2
}

variable "service_max_capacity" {
  description = "Maximum number of instances of the task for auto scaling"
  type        = number
  default     = 4
}

variable "log_retention_days" {
  description = "Number of days to retain CloudWatch Logs"
  type        = number
  default     = 30
} 