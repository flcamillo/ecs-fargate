# cria o endpoint para o serviço S3
resource "aws_vpc_endpoint" "endpoint_s3" {
  vpc_id            = aws_vpc.vpc.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = [for route in aws_route_table.route_table_private : route.id]
  tags = {
    Project = "${var.project_name}"
  }
}

# cria o endpoint para o serviço ecr - repositório
resource "aws_vpc_endpoint" "endpoint_ecr" {
  vpc_id              = aws_vpc.vpc.id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true
  security_group_ids  = [aws_security_group.sg_vpc_endpoint.id]
  subnet_ids          = [for subnet in aws_subnet.private_subnet : subnet.id]
  tags = {
    Project = "${var.project_name}"
  }
}

# cria o endpoint para o serviço ecr - api
resource "aws_vpc_endpoint" "endpoint_ecr_api" {
  vpc_id              = aws_vpc.vpc.id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.ecr.api"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true
  security_group_ids  = [aws_security_group.sg_vpc_endpoint.id]
  subnet_ids          = [for subnet in aws_subnet.private_subnet : subnet.id]
  tags = {
    Project = "${var.project_name}"
  }
}

# cria o endpoint para o serviço do cloudwatch
resource "aws_vpc_endpoint" "endpoint_cloudwatch" {
  vpc_id              = aws_vpc.vpc.id
  service_name        = "com.amazonaws.${data.aws_region.current.name}.logs"
  vpc_endpoint_type   = "Interface"
  private_dns_enabled = true
  security_group_ids  = [aws_security_group.sg_vpc_endpoint.id]
  subnet_ids          = [for subnet in aws_subnet.private_subnet : subnet.id]
  tags = {
    Project = "${var.project_name}"
  }
}
