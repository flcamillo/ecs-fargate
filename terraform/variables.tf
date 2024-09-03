variable "ecs_cluster_name" {
  description = "Nome do Cluster ECS"
  type        = string
  default     = "ECS-Fargate"
}

variable "ecr_repository_name" {
  description = "Nome do Repositório ECR para a Aplicação"
  type        = string
  default     = "go-app"
}

variable "ecr_repository_about" {
  description = "Descrição do Repositório ECR para a Aplicação"
  type        = string
  default     = "# GO Consumer SQS"
}

variable "ecr_repository_description" {
  description = "Descrição do Repositório ECR para a Aplicação"
  type        = string
  default     = "Aplicação em GO para consumir mensagens de fila SQS"
}

variable "service_name" {
  description = "Nome do serviço que irá subir no container"
  type        = string
  default     = "SQSConsumer"
}

variable "bucket_name" {
  description = "Nome do bucket S3 onde serão gravados os arquivos"
  type        = string
  default     = "files"
}

variable "sns_name" {
  description = "Nome do tópico SNS para notificar arquivos criados no bucket"
  type        = string
  default     = "notify-s3-object-created"
}

variable "sqs_name" {
  description = "Nome do SQS para notificar a lambda sobre arquivos criados no bucket, ela receberá eventos do SNS"
  type        = string
  default     = "process-s3-object-created"
}
