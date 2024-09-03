# identifica a conta atual
data "aws_caller_identity" "current" {}

# identifica a org da conta
data "aws_organizations_organization" "current" {}

# identifica a vpc da conta
data "aws_vpc" "current" {}

# region
data "aws_region" "current" {}

# identifica as subnets da vpc
data "aws_subnets" "current" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.current.id]
  }
}
