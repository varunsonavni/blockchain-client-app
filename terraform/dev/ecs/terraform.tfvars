# Environment settings
environment = "dev"
region      = "us-east-1"

container_version  = "latest"
blockchain_rpc_url = "https://polygon-rpc.com/"

# Task configuration
task_cpu    = 512  # 0.5 vCPU
task_memory = 1024 # 1 GB

# Service configuration
service_desired_count = 2
service_min_capacity  = 2
service_max_capacity  = 4

# Log settings
log_retention_days = 7
