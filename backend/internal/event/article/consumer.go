package article

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"github.com/johnwongx/webook/backend/pkg/saramax"
	"time"
)

type Consumer interface {
	Start() error
}

type KafkaConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.Logger
}

func NewKafkaConsumer(client sarama.Client, repo repository.InteractiveRepository, l logger.Logger) Consumer {
	return &KafkaConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (k *KafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", k.client)
	if err != nil {
		return err
	}

	go func() {
		err := cg.Consume(context.Background(),
			[]string{ReadEventTopic},
			saramax.NewConsumerHandler[ReadEvent](k.Consume, k.l))
		if err != nil {
			k.l.Error("消费循环退出异常", logger.Error(err))
		}
	}()
	return nil
}

func (k *KafkaConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return k.repo.IncrLike(ctx, evt.Aid, evt.Biz, evt.Uid)
}
