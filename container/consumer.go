package main

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Define o consumidor de mensagens da fila SQS que irá receber
// mensagens da fila SQS e envia-las para um canal para serem
// processadas.
type SQSConsumer struct {
	// Indica que o processo do consumidor foi encerrado
	stopped bool
	// Define um canal para aguardar o final do encerramento dos processos
	cancel chan bool
	// Canal para enviar as mensagens recebidas do SQS
	messages chan<- *types.Message
	// URL da fila SQS para excluir as mensagens processadas
	queueURL string
	// Define o log para enviar as mensagens
	log *Log
	// Client SQS para remover as mensagens da fila
	sqsClient SQSClient
}

// Cria o consumidor de mensagens SQS no canal informado.
func NewSQSConsumer(messages chan<- *types.Message, queueURL string, logOutput io.Writer) *SQSConsumer {
	return &SQSConsumer{
		messages: messages,
		queueURL: queueURL,
		cancel:   make(chan bool),
		log: &Log{
			o: logOutput,
		},
	}
}

// Inicializa o processo para receber as mensagens da fila SQS.
// Cada mensagem recebida será enviada para o canal para serem processadas.
func (p *SQSConsumer) Start() error {
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	if p.sqsClient == nil {
		p.sqsClient = sqs.NewFromConfig(sdkConfig)
	}
	for !p.stopped {
		received, err := p.sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            &p.queueURL,
			MaxNumberOfMessages: 10,
			VisibilityTimeout:   60,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			p.log.Error("failed to read queue, %s", err)
			return err
		}
		for _, message := range received.Messages {
			if p.stopped {
				break
			}
			p.messages <- &message
		}
	}
	p.cancel <- true
	return nil
}

// Encerra o consumo de mensagens
func (p *SQSConsumer) Stop() {
	p.stopped = true
	<-p.cancel
}
