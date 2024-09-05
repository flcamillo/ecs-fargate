# define o serviço para ativar no ecs
resource "aws_ecs_service" "ecs_service" {
  name            = var.service_name
  cluster         = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.task_definition.arn
  desired_count   = 1
  network_configuration {
    subnets          = [for subnet in aws_subnet.private_subnet : subnet.id]
    assign_public_ip = false
  }
  tags = {
    Project = "${var.project_name}"
  }
}
