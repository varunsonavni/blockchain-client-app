module "vpc" {
  source              = "../../modules/vpc"
  vpc_name            = var.vpc_name
  vpc_cidr            = var.vpc_cidr
  region              = var.region
  vpc_private_subnets = var.vpc_private_subnets
  vpc_public_subnets  = var.vpc_public_subnets
  environment         = var.environment

  public_subnet_tags = {
    "kubernetes.io/role/elb" : 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" : 1
  }
}