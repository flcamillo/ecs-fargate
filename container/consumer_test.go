package main

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Define a estrutura da mensagem gerada pelo SNS.
// A mensagem de criação do objeto do bucket S3 será adicionado no campo "Message".
var snsConsumerEventModel = `{
  "Type": "Notification",
  "MessageId": "dc1e94d9-56c5-5e96-808d-cc7f68faa162",
  "TopicArn": "arn:aws:sns:us-east-2:111122223333:ExampleTopic1",
  "Subject": "TestSubject",
  "Message": "%s",
  "Timestamp": "2021-02-16T21:41:19.978Z",
  "SignatureVersion": "1",
  "Signature": "FMG5tlZhJNHLHUXvZgtZzlk24FzVa7oX0T4P03neeXw8ZEXZx6z35j2FOTuNYShn2h0bKNC/zLTnMyIxEzmi2X1shOBWsJHkrW2xkR58ABZF+4uWHEE73yDVR4SyYAikP9jstZzDRm+bcVs8+T0yaLiEGLrIIIL4esi1llhIkgErCuy5btPcWXBdio2fpCRD5x9oR6gmE/rd5O7lX1c1uvnv4r1Lkk4pqP2/iUfxFZva1xLSRvgyfm6D9hNklVyPfy+7TalMD0lzmJuOrExtnSIbZew3foxgx8GT+lbZkLd0ZdtdRJlIyPRP44eyq78sU0Eo/LsDr0Iak4ZDpg8dXg==",
  "SigningCertURL": "https://sns.us-east-2.amazonaws.com/SimpleNotificationService-010a507c1833636cd94bdb98bd93083a.pem",
  "UnsubscribeURL": "https://sns.us-east-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-2:111122223333:ExampleTopic1:e1039402-24e7-40a3-a0d4-797da162b297"
}`

// Define a estrutura da mensagem gerada pelo S3.
// O nome do bucket e o nome arquivo serão adicionados na mensagem nos campos "name" e "key"
var s3ConsumerEventModel = `{  
   "Records":[  
      {  
         "eventVersion":"2.1",
         "eventSource":"aws:s3",
         "awsRegion":"us-west-2",
         "eventTime":"1970-01-01T00:00:00.000Z",
         "eventName":"ObjectCreated:Put",
         "userIdentity":{  
            "principalId":"AIDAJDPLRKLG7UEXAMPLE"
         },
         "requestParameters":{  
            "sourceIPAddress":"127.0.0.1"
         },
         "responseElements":{  
            "x-amz-request-id":"C3D13FE58DE4C810",
            "x-amz-id-2":"FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD"
         },
         "s3":{  
            "s3SchemaVersion":"1.0",
            "configurationId":"testConfigRule",
            "bucket":{  
               "name":"%s",
               "ownerIdentity":{  
                  "principalId":"A3NL1KOZZKExample"
               },
               "arn":"arn:aws:s3:::mybucket"
            },
            "object":{  
               "key":"%s",
               "size":1024,
               "eTag":"d41d8cd98f00b204e9800998ecf8427e",
               "versionId":"096fKKXTRTtl3on89fVO.nfljtsv6qko",
               "sequencer":"0055AED6DCD90281E5"
            }
         }
      }
   ]
}`

// Define uma estrutura para o SQS Client para usar como mock
type mockSQSConsumerClient struct{}

// Implementa as função de recebimento de mensagens simulando uma mensagem de retorno
func (p *mockSQSConsumerClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	s3Event := fmt.Sprintf(s3ConsumerEventModel, "bucket-teste", "teste.txt")
	s3Event = strings.ReplaceAll(s3Event, "\n", "")
	s3Event = strings.ReplaceAll(s3Event, `"`, `\"`)
	snsEvent := fmt.Sprintf(snsConsumerEventModel, s3Event)
	return &sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{
				Body:          &snsEvent,
				MessageId:     aws.String("123"),
				ReceiptHandle: aws.String("456"),
			},
		},
	}, nil
}

// Implementa a função de exclusão de mensagens, não gerando erro para simular um sucesso
func (p *mockSQSConsumerClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}

func TestConsumer(t *testing.T) {
	messages := make(chan *types.Message, 10)
	buffer := &bytes.Buffer{}
	consumer := NewSQSConsumer(messages, "", buffer)
	consumer.sqsClient = &mockSQSConsumerClient{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := consumer.Start()
		if err != nil {
			t.Error(err)
		}
		wg.Done()
	}()
	message := <-messages
	if *message.MessageId != "123" {
		t.Fatal("message id invalid")
	}
	if *message.ReceiptHandle != "456" {
		t.Fatal("message handle invalid")
	}
	consumer.Stop()
	wg.Wait()
}
