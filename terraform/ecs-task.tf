# define o padrão de tasks que serão executadas
resource "aws_ecs_task_definition" "task_definition" {
  family                   = var.service_name
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.role_task_execution.arn
  task_role_arn            = aws_iam_role.role_task.arn
  runtime_platform {
    cpu_architecture        = "X86_64"
    operating_system_family = "LINUX"
  }
  container_definitions = jsonencode([
    {
      name      = "${var.service_name}"
      image     = "900155862302.dkr.ecr.sa-east-1.amazonaws.com/go-app:1.0"
      cpu       = 256
      memory    = 512
      essential = true
      environment = [
        {
          name  = "TARGET_BUCKET"
          value = "${aws_s3_bucket.bucket_final.id}"
        },
        {
          name  = "QUEUE_URL"
          value = "${aws_sqs_queue.sqs.url}"
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-create-group  = "true"
          awslogs-group         = "${aws_cloudwatch_log_group.task_definition_log_group.name}"
          awslogs-region        = "${data.aws_region.current.name}"
          awslogs-stream-prefix = "/ecs/${var.ecs_cluster_name}"
        }
      }
    }
  ])
  tags = {
    Project = "${var.project_name}"
  }
}

# define o log group para a tarefa
resource "aws_cloudwatch_log_group" "task_definition_log_group" {
  name              = "/ecs/${var.ecs_cluster_name}/${var.service_name}"
  retention_in_days = 7
  tags = {
    Project = "${var.project_name}"
  }
}

# define a política para a tarefa
data "aws_iam_policy_document" "task_execution_policy_doc" {
  statement {
    sid       = "AllowAccessForLogs"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
    ]
  }
  statement {
    sid       = "AllowContainerInsights"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "events:PutRule",
      "events:PutTargets",
      "events:DescribeRule",
      "events:ListTargetsByRule",
    ]
  }
  statement {
    sid       = "AllowECRAccess"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:BatchGetImage",
      "ecr:DescribeImages",
      "ecr:DescribeImageScanFindings",
      "ecr:DescribeRepositories",
      "ecr:GetAuthorizationToken",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetLifecyclePolicy",
      "ecr:GetLifecyclePolicyPreview",
      "ecr:GetRepositoryPolicy",
      "ecr:ListImages",
      "ecr:ListTagsForResource"
    ]
  }
  statement {
    sid       = "AllowKMSAccess"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "kms:Decrypt",
    ]
  }
  statement {
    sid       = "AllowSecretManagerAccess"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "secretsmanager:GetSecretValue",
    ]
  }
  statement {
    sid       = "AllowSSM"
    effect    = "Allow"
    resources = ["*"]
    actions = [
      "ssm:GetParameters",
    ]
  }
}

# cria a policy da lambda com os acessos definidos
resource "aws_iam_policy" "task_execution_policy" {
  name        = "${var.ecs_cluster_name}_task_execution_policy"
  policy      = data.aws_iam_policy_document.task_execution_policy_doc.json
  path        = "/"
  description = "Policy para a Task Definition do ECS"
  tags = {
    Project = "${var.project_name}"
  }
}

# define a trusted policy da tarefa
data "aws_iam_policy_document" "task_definition_trusted_policy_doc" {
  statement {
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = [
        "ecs-tasks.amazonaws.com",
      ]
    }
    actions = ["sts:AssumeRole"]
  }
}

# cria a role da tarefa com a trusted policy
resource "aws_iam_role" "role_task_execution" {
  name               = "${var.ecs_cluster_name}_task_execution_role"
  assume_role_policy = data.aws_iam_policy_document.task_definition_trusted_policy_doc.json
  tags = {
    Project = "${var.project_name}"
  }
}

# associa a role com a policy
resource "aws_iam_role_policy_attachment" "attach_role_task_execution" {
  role       = aws_iam_role.role_task_execution.name
  policy_arn = aws_iam_policy.task_execution_policy.arn
}

# define a política para a tarefa
data "aws_iam_policy_document" "task_policy_doc" {
  statement {
    sid       = "SQSConsumeMessages"
    effect    = "Allow"
    resources = [aws_sqs_queue.sqs.arn]
    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
    ]
  }
  statement {
    sid    = "S3AccessSource"
    effect = "Allow"
    resources = [
      aws_s3_bucket.bucket_source.arn,
      "${aws_s3_bucket.bucket_source.arn}/*",
    ]
    actions = [
      "s3:List*",
      "s3:Get*",
      "s3:DeleteObject",
    ]
  }
  statement {
    sid    = "S3AccessFinal"
    effect = "Allow"
    resources = [
      aws_s3_bucket.bucket_final.arn,
      "${aws_s3_bucket.bucket_final.arn}/*",
    ]
    actions = [
      "s3:Put*",
    ]
  }
}

# cria a policy da lambda com os acessos definidos
resource "aws_iam_policy" "task_policy" {
  name        = "${var.ecs_cluster_name}_task_policy"
  policy      = data.aws_iam_policy_document.task_policy_doc.json
  path        = "/"
  description = "Policy para a Task Definition do ECS"
  tags = {
    Project = "${var.project_name}"
  }
}

# define a trusted policy da tarefa
data "aws_iam_policy_document" "task_trusted_policy_doc" {
  statement {
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = [
        "ecs-tasks.amazonaws.com",
      ]
    }
    actions = ["sts:AssumeRole"]
  }
}

# cria a role da tarefa com a trusted policy
resource "aws_iam_role" "role_task" {
  name               = "${var.ecs_cluster_name}_task_role"
  assume_role_policy = data.aws_iam_policy_document.task_trusted_policy_doc.json
  tags = {
    Project = "${var.project_name}"
  }
}

# associa a role com a policy
resource "aws_iam_role_policy_attachment" "attach_role_task" {
  role       = aws_iam_role.role_task.name
  policy_arn = aws_iam_policy.task_policy.arn
}
