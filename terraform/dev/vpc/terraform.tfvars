vpc_name            = "mm-test"
region              = "us-east-1"
vpc_cidr            = "10.0.0.0/16"
environment         = "dev"
vpc_private_subnets = ["10.0.0.0/20", "10.0.16.0/20"]
vpc_public_subnets  = ["10.0.128.0/20", "10.0.144.0/20"]
