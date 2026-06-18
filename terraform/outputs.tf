output "vpc_id" {
  value       = aws_vpc.main.id
  description = "VPC ID"
}

output "postgres_endpoint" {
  value       = aws_db_instance.postgres.endpoint
  description = "Connection endpoint for the RDS PostgreSQL database"
}

output "redis_endpoint" {
  value       = aws_elasticache_serverless_cache.redis.endpoint[0].address
  description = "Connection endpoint address for Redis Cache"
}

output "redis_port" {
  value       = aws_elasticache_serverless_cache.redis.endpoint[0].port
  description = "Connection endpoint port for Redis Cache"
}

output "rabbitmq_amqp_endpoints" {
  value       = aws_mq_broker.rabbitmq.instances[0].endpoints[0]
  description = "AMQP endpoint for RabbitMQ Broker connections"
}

output "rabbitmq_console_url" {
  value       = aws_mq_broker.rabbitmq.instances[0].console_url
  description = "Console URL for RabbitMQ Broker administration"
}

output "s3_user_avatars_bucket" {
  value       = aws_s3_bucket.user_avatars.id
  description = "S3 Bucket Name for user avatars"
}

output "s3_blog_covers_bucket" {
  value       = aws_s3_bucket.blog_covers.id
  description = "S3 Bucket Name for blog covers"
}
