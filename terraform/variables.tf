variable "aws_region" {
  type        = string
  description = "The AWS region to deploy resources into"
  default     = "ap-south-1"
}

variable "project_name" {
  type        = string
  description = "The name of the project"
  default     = "blog-platform"
}

variable "environment" {
  type        = string
  description = "Deployment environment"
  default     = "production"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
  default     = "10.0.0.0/16"
}

variable "db_username" {
  type        = string
  description = "Username for the RDS PostgreSQL database"
  default     = "postgres"
}

variable "db_password" {
  type        = string
  description = "Password for the RDS PostgreSQL database"
  sensitive   = true
}

variable "redis_node_type" {
  type        = string
  description = "Node size for ElastiCache Redis"
  default     = "cache.t4g.micro"
}

variable "rabbitmq_user" {
  type        = string
  description = "Username for Amazon MQ RabbitMQ broker"
  default     = "admin"
}

variable "rabbitmq_password" {
  type        = string
  description = "Password for Amazon MQ RabbitMQ broker"
  sensitive   = true
}
