package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Define o processador de mensagens da fila SQS
type Worker struct {
	// Canal para receber as mensagens
	messages <-chan *types.Message
	// URL da fila SQS para excluir as mensagens processadas
	queueURL string
	// Bucket para copiar o arquivo
	targetBucket string
	// Define o log para enviar as mensagens
	log *Log
	// Client S3 para realizar a copia entre os buckets
	s3Client S3Client
	// Client SQS para remover as mensagens da fila
	sqsClient SQSClient
}

// Cria o processador de mensagens no canal informado.
func NewWorker(messages <-chan *types.Message, queueURL string, targetBucket string, logOutput io.Writer) *Worker {
	return &Worker{
		messages:     messages,
		queueURL:     queueURL,
		targetBucket: targetBucket,
		log: &Log{
			o: logOutput,
		},
	}
}

// Inicializa o processamento das mensagens recebidas no canal.
func (p *Worker) Start() error {
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	if p.s3Client == nil {
		p.s3Client = s3.NewFromConfig(sdkConfig)
	}
	if p.sqsClient == nil {
		p.sqsClient = sqs.NewFromConfig(sdkConfig)
	}
	for {
		message, ok := <-p.messages
		if !ok {
			break
		}
		p.log.SetCorrelation(*message.MessageId)
		err := p.processMessage(message)
		if err != nil {
			p.log.Error("failed to process message, %s", err)
		}
		p.removeMessage(message)
	}
	return nil
}

// Processa a mensagem extraindo a mensagem SNS do corpo da mensagem SQS e
// em seguida extrai da mensagem SNS os eventos de gravação no bucket S3.
// Para cada evento S3 recebido executa o processo de copia do arquivo
// do bucket original para o bucket de destino e remove o arquivo após a
// cópia com sucesso.
func (p *Worker) processMessage(message *types.Message) error {
	snsEvent := &events.SNSEntity{}
	err := json.Unmarshal([]byte(*message.Body), &snsEvent)
	if err != nil {
		return fmt.Errorf("failed to decode sns event, %s", err)
	}
	s3Event := &events.S3Event{}
	err = json.Unmarshal([]byte(snsEvent.Message), &s3Event)
	if err != nil {
		return fmt.Errorf("failed to decode s3 event, %s", err)
	}
	received := len(s3Event.Records)
	if received < 1 {
		return nil
	}
	for _, record := range s3Event.Records {
		start := time.Now()
		p.log.Info("starting copy of {%s/%s} to {%s/%s}...", record.S3.Bucket.Name, record.S3.Object.Key, p.targetBucket, record.S3.Object.Key)
		err = p.processCopy(&record, true)
		elapsed := time.Since(start)
		if err != nil {
			p.log.Error("copy of {%s/%s} to {%s/%s} failed, %s", record.S3.Bucket.Name, record.S3.Object.Key, p.targetBucket, record.S3.Object.Key, err)
			continue
		}
		rate := float64(record.S3.Object.Size)
		if elapsed.Seconds() > 0 {
			rate = float64(record.S3.Object.Size) / elapsed.Seconds()
		}
		p.log.Info("copy of {%s/%s} to {%s/%s} completed successfully in {%d ms} rate {%.2f bytes/s}", record.S3.Bucket.Name, record.S3.Object.Key, p.targetBucket, record.S3.Object.Key, elapsed.Microseconds(), rate)
	}
	return nil
}

// Executa a cópia do arquivo do bucket original para o destino e faz a
// exclusão do arquivo original após cópia com sucesso se solicitado.
func (p *Worker) processCopy(event *events.S3EventRecord, remove bool) error {
	sourceFile := fmt.Sprintf("%s/%s", event.S3.Bucket.Name, event.S3.Object.Key)
	_, err := p.s3Client.CopyObject(
		context.Background(),
		&s3.CopyObjectInput{
			Bucket:     &p.targetBucket,
			Key:        &event.S3.Object.Key,
			CopySource: &sourceFile,
		})
	if err != nil {
		return err
	}
	if remove {
		_, err = p.s3Client.DeleteObject(
			context.Background(),
			&s3.DeleteObjectInput{
				Bucket: &p.targetBucket,
				Key:    &event.S3.Object.Key,
			})
		if err != nil {
			p.log.Warning("failed to delete source file, %s", err)
		}
	}
	return nil
}

// Remove a mensagem da fila SQS para evitar processá-la novamente.
func (p *Worker) removeMessage(message *types.Message) {
	_, err := p.sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      &p.queueURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		p.log.Warning("unable to delete message from sqs, %s", err)
	}
}
