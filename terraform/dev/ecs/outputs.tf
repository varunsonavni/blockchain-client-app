
output "ecr_repository_url" {
  description = "The URL of the ECR Repository"
  value       = aws_ecr_repository.blockchain_client.repository_url
}

output "ecs_cluster_name" {
  description = "The name of the ECS cluster"
  value       = module.ecs.cluster_name
}

output "ecs_service_name" {
  description = "The name of the ECS service"
  value       = module.ecs.service_name
}

output "load_balancer_dns" {
  description = "The DNS name of the load balancer"
  value       = module.ecs.alb_dns_name
}

output "blockchain_client_url" {
  description = "URL to access the blockchain client API"
  value       = "http://${module.ecs.alb_dns_name}"
} 