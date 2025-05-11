module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.vpc_name}-${var.environment}"
  cidr = var.vpc_cidr

  azs             = ["${var.region}a", "${var.region}c"]
  private_subnets = var.vpc_private_subnets
  public_subnets  = var.vpc_public_subnets

  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_nat_gateway  = true
  single_nat_gateway  = true
  reuse_nat_ips       = true # <= Skip creation of EIPs for the NAT Gateways
  external_nat_ip_ids = aws_eip.nat.*.id

  manage_default_network_acl = true

  manage_default_route_table = true

  manage_default_security_group = true

  public_subnet_tags  = var.public_subnet_tags
  private_subnet_tags = var.private_subnet_tags

  tags = {
    Environment = "${var.environment}"
  }
}

resource "aws_eip" "nat" {
  count = 1
}

