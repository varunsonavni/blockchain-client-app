variable "environment" {
  description = "Deployment environment"
  type        = string
}

variable "region" {
  description = "AWS region"
  type        = string
}

variable "public_subnet_tags" {
  description = "Tags for public subnets"
  type        = map(string)
  default = {
    "Tier" = "Public"
  }
}

variable "private_subnet_tags" {
  description = "Tags for private subnets"
  type        = map(string)
  default = {
    "Tier" = "Private"
  }
}

variable "container_version" {
  description = "Docker image version/tag for the blockchain client"
  type        = string
  default     = "latest"
}

variable "blockchain_rpc_url" {
  description = "Blockchain RPC URL (e.g., Polygon RPC endpoint)"
  type        = string
}

variable "task_cpu" {
  description = "CPU units for the ECS task (1024 = 1 vCPU)"
  type        = number
}

variable "task_memory" {
  description = "Memory for the ECS task in MiB"
  type        = number
}

variable "service_desired_count" {
  description = "Number of instances of the task to run"
  type        = number
}

variable "service_min_capacity" {
  description = "Minimum number of instances of the task for auto scaling"
  type        = number
}

variable "service_max_capacity" {
  description = "Maximum number of instances of the task for auto scaling"
  type        = number
}

variable "log_retention_days" {
  description = "Number of days to retain CloudWatch Logs"
  type        = number
} 