package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const ReadEventTopic = "topic_read_event"

type Producer interface {
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: producer,
	}
}

func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, evt ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: ReadEventTopic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

type ReadEvent struct {
	Uid int64
	Aid int64
	Biz string
}
