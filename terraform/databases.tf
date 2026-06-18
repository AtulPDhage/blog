# Subnet group for RDS Databases
resource "aws_db_subnet_group" "rds" {
  name       = "${var.project_name}-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  tags = {
    Name = "${var.project_name}-db-subnet-group"
  }
}

# RDS PostgreSQL Database Instance
resource "aws_db_instance" "postgres" {
  identifier             = "${var.project_name}-postgres"
  engine                 = "postgres"
  engine_version         = "15.4"
  instance_class         = "db.t4g.micro" # Lightweight cost-effective instance size
  allocated_storage      = 20
  max_allocated_storage  = 100
  db_name                = "postgres"
  username               = var.db_username
  password               = var.db_password
  db_subnet_group_name   = aws_db_subnet_group.rds.name
  vpc_security_group_ids = [aws_security_group.db.id]
  skip_final_snapshot    = true

  tags = {
    Name        = "${var.project_name}-db"
    Environment = var.environment
  }
}

# Subnet group for Redis Cache
resource "aws_elasticache_subnet_group" "redis" {
  name       = "${var.project_name}-redis-subnet-group"
  subnet_ids = aws_subnet.private[*].id
}

# Redis Serverless Cache Cluster
# Serverless Redis is highly scalable and automatically scales to handle spikes in traffic
resource "aws_elasticache_serverless_cache" "redis" {
  name       = "${var.project_name}-redis"
  engine     = "redis"
  major_engine_version = "7"
  
  subnet_ids = aws_subnet.private[*].id
  security_group_ids = [aws_security_group.redis.id]

  tags = {
    Name        = "${var.project_name}-redis"
    Environment = var.environment
  }
}
