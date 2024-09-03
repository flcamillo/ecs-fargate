# define o serviço para ativar no ecs
resource "aws_ecs_service" "ecs_service" {
  name            = var.service_name
  cluster         = aws_ecs_cluster.ecs_cluster.id
  task_definition = aws_ecs_task_definition.task_definition.arn
  desired_count   = 1
  network_configuration {
    subnets          = data.aws_subnets.current.ids
    assign_public_ip = "true"
    #security_groups = [aws_security_group.ecs_service_security_group.id]
  }
}

# # define o grupo de segurança do serviço
# resource "aws_security_group" "ecs_service_security_group" {
#   vpc_id      = data.aws_vpc.current.id
#   name        = "AllowTCP"
#   description = "Permite trafego TCP IP para o container e qualquer saída"
# }

# # define a regra de entrada no grupo de segurança
# resource "aws_vpc_security_group_ingress_rule" "ingress_ecs_service" {
#   security_group_id = aws_security_group.ecs_service_security_group.id
#   cidr_ipv4         = data.aws_vpc.current.cidr_block
#   from_port         = 80
#   ip_protocol       = "tcp"
#   to_port           = 80
# }

# # define a regra de sa~ida no grupo de segurança
# resource "aws_vpc_security_group_egress_rule" "egress_ecs_service" {
#   security_group_id = aws_security_group.ecs_service_security_group.id
#   cidr_ipv4         = "0.0.0.0/0"
#   ip_protocol       = "-1"
# }
