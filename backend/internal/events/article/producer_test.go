//go:build e2e

package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

func TestProducer(t *testing.T) {
	client := InitKafka()
	syncProducer := NewSyncProducer(client)
	producer := NewKafkaProducer(syncProducer)
	err := producer.ProduceReadEvent(context.Background(), ReadEvent{
		Uid: 123,
		Aid: 1,
		Biz: "article",
	})
	assert.NoError(t, err)
}
