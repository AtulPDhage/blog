# Amazon MQ RabbitMQ Broker Configuration
resource "aws_mq_broker" "rabbitmq" {
  broker_name = "${var.project_name}-rabbitmq"

  engine_type        = "RabbitMQ"
  engine_version     = "3.10.8"
  host_instance_type = "mq.t3.micro" # Lightweight managed node
  deployment_mode    = "SINGLE_INSTANCE"

  user {
    username = var.rabbitmq_user
    password = var.rabbitmq_password
  }

  subnet_ids         = [aws_subnet.private[0].id] # Provisioned inside our private subnet for security
  security_groups    = [aws_security_group.mq.id]
  publicly_accessible = false # Keeping it private restricts queue access to internal microservices

  tags = {
    Name        = "${var.project_name}-rabbitmq"
    Environment = var.environment
  }
}
