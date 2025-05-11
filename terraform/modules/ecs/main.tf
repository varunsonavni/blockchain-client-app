resource "aws_ecs_cluster" "main" {
  name = "${var.environment}-blockchain-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled"
  }

  tags = {
    Name        = "${var.environment}-blockchain-cluster"
    Environment = var.environment
  }
}

# Security group for ALB
resource "aws_security_group" "alb" {
  name        = "${var.environment}-alb-sg"
  description = "Security group for blockchain client ALB"
  vpc_id      = var.vpc_id

  # Allow HTTP from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow HTTPS from anywhere
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.environment}-alb-sg"
    Environment = var.environment
  }
}

# Security group for ECS tasks
resource "aws_security_group" "ecs_tasks" {
  name        = "${var.environment}-ecs-tasks-sg"
  description = "Security group for blockchain client ECS tasks"
  vpc_id      = var.vpc_id

  # Allow inbound from ALB
  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.environment}-ecs-tasks-sg"
    Environment = var.environment
  }
}

# ECS Task Execution Role
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.environment}-ecs-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

# Attach the ECS Task Execution policy to the role
resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Task Definition for the blockchain client
resource "aws_ecs_task_definition" "blockchain_client" {
  family                   = "${var.environment}-blockchain-client"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "blockchain-client"
      image     = "${var.container_image}:${var.container_version}"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "BLOCKCHAIN_RPC_URL"
          value = var.blockchain_rpc_url
        },
        {
          name  = "API_PORT"
          value = ":8080"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${var.environment}-blockchain-client"
          "awslogs-region"        = var.region
          "awslogs-stream-prefix" = "ecs"
        }
      }
      healthCheck = {
        command     = ["CMD-SHELL", "wget -q --spider http://localhost:8080/api/blocks/latest || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }
    }
  ])

  tags = {
    Name        = "${var.environment}-blockchain-client"
    Environment = var.environment
  }
}

# CloudWatch Log Group for ECS logs
resource "aws_cloudwatch_log_group" "blockchain_client" {
  name              = "/ecs/${var.environment}-blockchain-client"
  retention_in_days = var.log_retention_days

  tags = {
    Name        = "${var.environment}-blockchain-client-logs"
    Environment = var.environment
  }
}

# Application Load Balancer
resource "aws_lb" "blockchain_client" {
  name               = "${var.environment}-blockchain-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = var.environment == "prod" ? true : false

  tags = {
    Name        = "${var.environment}-blockchain-alb"
    Environment = var.environment
  }
}

# ALB Target Group
resource "aws_lb_target_group" "blockchain_client" {
  name        = "${var.environment}-blockchain-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    path                = "/api/blocks/latest"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    protocol            = "HTTP"
    matcher             = "200"
  }

  tags = {
    Name        = "${var.environment}-blockchain-tg"
    Environment = var.environment
  }
}

# ALB Listener
resource "aws_lb_listener" "blockchain_client" {
  load_balancer_arn = aws_lb.blockchain_client.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.blockchain_client.arn
  }
}

# ECS Service
resource "aws_ecs_service" "blockchain_client" {
  name                               = "${var.environment}-blockchain-service"
  cluster                            = aws_ecs_cluster.main.id
  task_definition                    = aws_ecs_task_definition.blockchain_client.arn
  desired_count                      = var.service_desired_count
  deployment_minimum_healthy_percent = 50
  deployment_maximum_percent         = 200
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"
  health_check_grace_period_seconds  = 60

  network_configuration {
    security_groups  = [aws_security_group.ecs_tasks.id]
    subnets          = var.private_subnet_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.blockchain_client.arn
    container_name   = "blockchain-client"
    container_port   = 8080
  }

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }

  depends_on = [aws_lb_listener.blockchain_client]

  tags = {
    Name        = "${var.environment}-blockchain-service"
    Environment = var.environment
  }
}

# Auto Scaling for ECS Service
resource "aws_appautoscaling_target" "blockchain_client" {
  max_capacity       = var.service_max_capacity
  min_capacity       = var.service_min_capacity
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.blockchain_client.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

# CPU Utilization Scaling Policy
resource "aws_appautoscaling_policy" "cpu" {
  name               = "${var.environment}-blockchain-cpu-autoscaling"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.blockchain_client.resource_id
  scalable_dimension = aws_appautoscaling_target.blockchain_client.scalable_dimension
  service_namespace  = aws_appautoscaling_target.blockchain_client.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }

    target_value       = 70.0
    scale_in_cooldown  = 300
    scale_out_cooldown = 300
  }
} 