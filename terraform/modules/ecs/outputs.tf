output "cluster_id" {
  description = "The ID of the ECS cluster"
  value       = aws_ecs_cluster.main.id
}

output "cluster_name" {
  description = "The name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "service_name" {
  description = "The name of the ECS service"
  value       = aws_ecs_service.blockchain_client.name
}

output "task_definition_arn" {
  description = "The ARN of the Task Definition"
  value       = aws_ecs_task_definition.blockchain_client.arn
}

output "alb_dns_name" {
  description = "The DNS name of the load balancer"
  value       = aws_lb.blockchain_client.dns_name
}

output "alb_zone_id" {
  description = "The canonical hosted zone ID of the load balancer"
  value       = aws_lb.blockchain_client.zone_id
}

output "log_group_name" {
  description = "The name of the CloudWatch log group"
  value       = aws_cloudwatch_log_group.blockchain_client.name
}

output "alb_security_group_id" {
  description = "The ID of the ALB security group"
  value       = aws_security_group.alb.id
}

output "ecs_security_group_id" {
  description = "The ID of the ECS tasks security group"
  value       = aws_security_group.ecs_tasks.id
} 