package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {
	// identifica o bucket de destino para os arquivos recebidos, a fila sqs
	// e o total de processos em paralelo
	targetBucket := os.Getenv("TARGET_BUCKET")
	queueURL := os.Getenv("QUEUE_URL")
	workers := 10
	if os.Getenv("WORKERS") != "" {
		n, err := strconv.Atoi(os.Getenv("WORKERS"))
		if err == nil && n > 0 {
			workers = n
		}
	}
	// inicializa o log
	logOutput := os.Stdout
	log := NewLog(logOutput)
	log.Info("using target bucket: %s", targetBucket)
	log.Info("using queue url: %s", queueURL)
	log.Info("using max workers: %d", workers)
	// cria o canal para processar as mensagens com o maximo de workers definido
	messages := make(chan *types.Message, workers)
	defer close(messages)
	// inicia os workers
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		worker := NewWorker(messages, targetBucket, queueURL, logOutput)
		wg.Add(1)
		go func() {
			err := worker.Start()
			if err != nil {
				log.Error("failed do start worker, %s", err)
			}
			wg.Done()
		}()
	}
	// inicia o consumidor de mensagens do SQS
	consumer := NewSQSConsumer(messages, queueURL, logOutput)
	wg.Add(1)
	go func() {
		err := consumer.Start()
		if err != nil {
			log.Error("failed do start consumer, %s", err)
		}
		wg.Done()
	}()
	// so encerra o program de todos os processos em background terminarem
	wg.Wait()
}
